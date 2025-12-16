package worker

import (
    "Gym_StrongCode/internal/repository"
    "Gym_StrongCode/internal/service"
    "context"
    "database/sql"
    "log"
    "time"
)

type NotificationWorker struct {
    notificationService *service.NotificationService
    userRepo           *repository.UserRepository
    bookingRepo        *repository.BookingRepository
    db                 *sql.DB
    interval           time.Duration
}

func NewNotificationWorker(notificationService *service.NotificationService, interval time.Duration) *NotificationWorker {
    return &NotificationWorker{
        notificationService: notificationService,
        interval:           interval,
    }
}

func (w *NotificationWorker) Start(ctx context.Context) {
    ticker := time.NewTicker(w.interval)
    defer ticker.Stop()

    log.Println("Notification worker started")
    
    for {
        select {
        case <-ctx.Done():
            log.Println("Notification worker stopped")
            return
        case <-ticker.C:
            w.processNotifications()
        }
    }
}

func (w *NotificationWorker) processNotifications() {
    // Отправляем напоминания о занятиях за 2 часа до начала
    w.sendClassReminders()
    
    // Проверяем истекшие подписки
    w.checkExpiredMemberships()
    
    // Очищаем старые уведомления
    w.cleanupOldNotifications()
}

func (w *NotificationWorker) sendClassReminders() {
    // Находим занятия, которые начнутся через 2 часа
    twoHoursFromNow := time.Now().Add(2 * time.Hour).Format("2006-01-02 15:04:05")
    
    // Здесь должен быть запрос к БД для получения занятий
    // и пользователей, которым нужно отправить напоминания
    // Примерная логика:
    /*
    rows, err := w.db.Query(`
        SELECT b.user_id, u.email, u.name, c.title, c.start_time
        FROM bookings b
        JOIN users u ON b.user_id = u.id
        JOIN classes c ON b.class_id = c.id
        WHERE b.status = 'booked'
        AND c.start_time BETWEEN datetime('now') AND ?
        AND b.reminder_sent = 0
    `, twoHoursFromNow)
    */
    
    log.Printf("Checking for class reminders (next 2 hours: %s)", twoHoursFromNow)
}

func (w *NotificationWorker) checkExpiredMemberships() {
    today := time.Now().Format("2006-01-02")
    
    // Находим подписки, которые истекают сегодня
    // Примерная логика:
    /*
    rows, err := w.db.Query(`
        SELECT um.user_id, u.email, u.name, m.name as membership_name, um.end_date
        FROM user_memberships um
        JOIN users u ON um.user_id = u.id
        JOIN memberships m ON um.membership_id = m.id
        WHERE um.active = 1
        AND um.end_date = ?
        AND um.expiry_notification_sent = 0
    `, today)
    */
    
    log.Printf("Checking for expiring memberships (today: %s)", today)
}

func (w *NotificationWorker) cleanupOldNotifications() {
    // Удаляем уведомления старше 30 дней
    thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
    
    // Примерная логика:
    /*
    _, err := w.db.Exec(`
        DELETE FROM notifications 
        WHERE created_at < ?
    `, thirtyDaysAgo)
    */
    
    log.Printf("Cleaning up old notifications (older than: %s)", thirtyDaysAgo)
}

func (w *NotificationWorker) sendEmailNotification(email, subject, body string) {
    err := w.notificationService.sendEmail(email, subject, body)
    if err != nil {
        log.Printf("Failed to send email to %s: %v", email, err)
    } else {
        log.Printf("Email sent to %s: %s", email, subject)
    }
}