package realtime

import (
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"ecochitas/internal/domain"
)

type truck_position_subscriber struct {
	subscriber_identifier   string
	truck_identifier_filter string
	truck_position_event_ch chan domain.Truck_latest_position
}

type Truck_position_stream struct {
	application_logger    *slog.Logger
	subscribers_by_id     map[string]truck_position_subscriber
	subscribers_mutex     sync.RWMutex
	next_subscriber_index atomic.Uint64
}

func New_truck_position_stream(application_logger *slog.Logger) *Truck_position_stream {
	return &Truck_position_stream{
		application_logger: application_logger,
		subscribers_by_id:  make(map[string]truck_position_subscriber),
	}
}

func (truck_position_stream *Truck_position_stream) Subscribe(
	raw_truck_identifier_filter string,
) (string, <-chan domain.Truck_latest_position, func()) {
	truck_identifier_filter := strings.TrimSpace(raw_truck_identifier_filter)
	subscriber_index := truck_position_stream.next_subscriber_index.Add(1)
	subscriber_identifier := "subscriber_" + format_uint64(subscriber_index)

	truck_position_event_ch := make(chan domain.Truck_latest_position, 64)
	subscriber_item := truck_position_subscriber{
		subscriber_identifier:   subscriber_identifier,
		truck_identifier_filter: truck_identifier_filter,
		truck_position_event_ch: truck_position_event_ch,
	}

	truck_position_stream.subscribers_mutex.Lock()
	truck_position_stream.subscribers_by_id[subscriber_identifier] = subscriber_item
	truck_position_stream.subscribers_mutex.Unlock()

	unsubscribe_function := func() {
		truck_position_stream.subscribers_mutex.Lock()
		subscriber_to_remove, has_subscriber := truck_position_stream.subscribers_by_id[subscriber_identifier]
		if has_subscriber {
			delete(truck_position_stream.subscribers_by_id, subscriber_identifier)
			close(subscriber_to_remove.truck_position_event_ch)
		}
		truck_position_stream.subscribers_mutex.Unlock()
	}

	return subscriber_identifier, truck_position_event_ch, unsubscribe_function
}

func (truck_position_stream *Truck_position_stream) Publish(truck_latest_position domain.Truck_latest_position) {
	truck_position_stream.subscribers_mutex.RLock()
	defer truck_position_stream.subscribers_mutex.RUnlock()

	for _, subscriber_item := range truck_position_stream.subscribers_by_id {
		if subscriber_item.truck_identifier_filter != "" &&
			subscriber_item.truck_identifier_filter != truck_latest_position.Truck_identifier {
			continue
		}

		select {
		case subscriber_item.truck_position_event_ch <- truck_latest_position:
		default:
			truck_position_stream.application_logger.Warn(
				"truck_position_stream_subscriber_is_slow_dropping_event",
				"subscriber_identifier",
				subscriber_item.subscriber_identifier,
				"truck_identifier",
				truck_latest_position.Truck_identifier,
			)
		}
	}
}

func format_uint64(integer_value uint64) string {
	const decimal_radix = 10
	return strconv.FormatUint(integer_value, decimal_radix)
}
