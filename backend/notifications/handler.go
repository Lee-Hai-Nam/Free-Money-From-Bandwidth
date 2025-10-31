package notifications

import (
	"fmt"
	"sync"
	"time"

	"bandwidth-income-manager/backend/monitor"
)

// Handler manages notifications and alerts
type Handler struct {
	config   *Config
	monitor  *monitor.Collector
	channels []NotificationChannel
	mu       sync.RWMutex
}

// NewHandler creates a new notification handler
func NewHandler(config *Config, monitor *monitor.Collector) *Handler {
	return &Handler{
		config:   config,
		monitor:  monitor,
		channels: make([]NotificationChannel, 0),
	}
}

// Config represents notification configuration
type Config struct {
	Enabled           bool
	AppStopped        bool
	EarningsMilestone bool
	UpdateAvailable   bool
	ProxyFailure      bool
	DiscordWebhook    string
	TelegramBotToken  string
	TelegramChatID    string
	EmailSMTPHost     string
	EmailSMTPPort     int
	EmailUsername     string
	EmailPassword     string
	EmailTo           string
}

// SendNotification sends a notification through all channels
func (h *Handler) SendNotification(event *NotificationEvent) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, channel := range h.channels {
		if err := channel.Send(event); err != nil {
			fmt.Printf("Failed to send notification via %s: %v\n", channel.Name(), err)
		}
	}

	return nil
}

// NotifyAppStopped notifies that an app has stopped
func (h *Handler) NotifyAppStopped(appID string) {
	if !h.config.AppStopped {
		return
	}

	event := &NotificationEvent{
		Type:      EventAppStopped,
		AppID:     appID,
		Message:   fmt.Sprintf("App %s has stopped running", appID),
		Timestamp: time.Now(),
	}

	h.SendNotification(event)
}

// NotifyEarningsMilestone notifies about earnings milestones
func (h *Handler) NotifyEarningsMilestone(amount float64, currency string) {
	if !h.config.EarningsMilestone {
		return
	}

	event := &NotificationEvent{
		Type:      EventEarningsMilestone,
		Message:   fmt.Sprintf("You've reached $%.2f %s in total earnings!", amount, currency),
		Timestamp: time.Now(),
		Metadata:  map[string]interface{}{"amount": amount, "currency": currency},
	}

	h.SendNotification(event)
}

// NotifyUpdateAvailable notifies about app updates
func (h *Handler) NotifyUpdateAvailable(appID string, version string) {
	if !h.config.UpdateAvailable {
		return
	}

	event := &NotificationEvent{
		Type:      EventUpdateAvailable,
		AppID:     appID,
		Message:   fmt.Sprintf("New version available for %s: %s", appID, version),
		Timestamp: time.Now(),
		Metadata:  map[string]interface{}{"version": version},
	}

	h.SendNotification(event)
}

// NotifyProxyFailure notifies about proxy failures
func (h *Handler) NotifyProxyFailure(proxyID string) {
	if !h.config.ProxyFailure {
		return
	}

	event := &NotificationEvent{
		Type:      EventProxyFailure,
		Message:   fmt.Sprintf("Proxy %s is not responding", proxyID),
		Timestamp: time.Now(),
		Metadata:  map[string]interface{}{"proxy_id": proxyID},
	}

	h.SendNotification(event)
}

// NotificationEvent represents a notification event
type NotificationEvent struct {
	Type      EventType
	AppID     string
	DeviceID  string
	Message   string
	Timestamp time.Time
	Metadata  map[string]interface{}
}

// EventType represents the type of notification event
type EventType string

const (
	EventAppStopped        EventType = "app_stopped"
	EventEarningsMilestone EventType = "earnings_milestone"
	EventUpdateAvailable   EventType = "update_available"
	EventProxyFailure      EventType = "proxy_failure"
)

// NotificationChannel interface for different notification channels
type NotificationChannel interface {
	Name() string
	Send(event *NotificationEvent) error
}

// NativeNotificationChannel sends native OS notifications
type NativeNotificationChannel struct{}

func (n *NativeNotificationChannel) Name() string {
	return "native"
}

func (n *NativeNotificationChannel) Send(event *NotificationEvent) error {
	// TODO: Implement native OS notifications
	return nil
}

// DiscordWebhookChannel sends notifications to Discord webhook
type DiscordWebhookChannel struct {
	WebhookURL string
}

func (d *DiscordWebhookChannel) Name() string {
	return "discord"
}

func (d *DiscordWebhookChannel) Send(event *NotificationEvent) error {
	// TODO: Implement Discord webhook
	return nil
}

// TelegramChannel sends notifications to Telegram
type TelegramChannel struct {
	BotToken string
	ChatID   string
}

func (t *TelegramChannel) Name() string {
	return "telegram"
}

func (t *TelegramChannel) Send(event *NotificationEvent) error {
	// TODO: Implement Telegram notifications
	return nil
}
