## Services:
1. Api-gateway
2. Enquiry-service -> Postgres DB
3. Mail-service
4. Logger-service -> Postgres DB
5. Listener-service
6. frontend

### DB tables:
- users table
- enquiry tables
- schedules table
- logs tables

### Microservice communication
- http
- gRPC
- Rabbitmq message broker

### Containerization
- Docker

### Databases
- Postgres -> data store
- Redis -> Cache, Rate Limiting

### 3rd party API Integration
- Sendgrid

### Service wise methods
> ### Api-gateway
##### HTTP
```go
func HandleSubmission(c *gin.Context)
func (lac *LocalApiConfig) FetchAllProperties(c *gin.Context)
func connect() (*amqp.Connection, error)
```
##### gRPC
```go
func EnquiryViaGRPC(c *gin.Context, enquiryPayload types.EnquiryPayload)
func (lac *LocalApiConfig) GetAllLogs(c *gin.Context)
func (lac *LocalApiConfig) WriteLog(c *gin.Context)
func (lac *LocalApiConfig) CreateNewUser(c *gin.Context, userPayload types.UserPayload)
```

> ### Enquiry-service
##### HTTP
```go
func (localApiConfig *LocalApiConfig) HandleFetchAllProperties(c *gin.Context)
```
```go
func StartGrpcServer(localApiConfig *handlers.LocalApiConfig)
func connect() (*amqp.Connection, error)
func connectToDB() *sql.DB
func ProcessScheduledTasks(s *handlers.EnquiryServer)
```
##### handlers
```go
func StartGrpcUserServer(localApiConfig *handlers.LocalApiConfig)
func (e *EnquiryServer) HandleCustomerEnquiry(ctx context.Context, request *enquiries.CustomerEnquiryRequest) (*enquiries.CustomerEnquiryResponse, error)
func (e *EnquiryServer) getTotalEnquiriesLastWeek(c context.Context, updatedUser database.User) (int, error)
func (e *EnquiryServer) SendEmail(payload EnquiryMailPayloadUsingSendgrid) error
func (e *EnquiryServer) notifyUserAboutEnquiry(input *enquiries.CustomerEnquiry, totalEnquiries int, mailPayload EnquiryMailPayloadUsingSendgrid) error
func (e *EnquiryServer) executeTask(input *enquiries.CustomerEnquiry, mailPayload EnquiryMailPayloadUsingSendgrid) error
func (e *EnquiryServer) scheduleTask(input *enquiries.CustomerEnquiry, mailPayload EnquiryMailPayloadUsingSendgrid) error
func (u *UserServer) CreateNewUser(ctx context.Context, request *users.CreateUserRequest) (*users.CreateUserResponse, error)
```
##### rabbitmq actions
```go
func (e *Emitter) declareExchange(ch *amqp.Channel) error
func (e *Emitter) Emit(event string) error
func NewEmitter(conn *amqp.Connection, exchange, routingKey string) (*Emitter, error)
```
> ### Listener-service
##### rabbitmq actions
```go
func DeclareExchange(ch *amqp.Channel) error
func DeclareRandomQueue(ch *amqp.Channel) (amqp.Queue, error)
func DeclareMailExchange(ch *amqp.Channel) error
func DeclareMailQueue(ch *amqp.Channel) (amqp.Queue, error)
func DeclareEnquiryMailQueue(ch *amqp.Channel) (amqp.Queue, error)
func NewLogConsumer(conn *amqp.Connection) (*LogConsumer, error)
func (consumer *LogConsumer) ConsumeLogs(topics []string) error
func logEvent(log LogPayload) error
func NewMailConsumer(conn *amqp.Connection) (*MailConsumer, error)
func (consumer *MailConsumer) ConsumeMails() error
func (consumer *MailConsumer) ConsumeEnquiryMails() error
func sendEnquiryMail(payload EnquiryMailPayloadUsingSendgrid) (time.Duration, error)
func logMailSendingResult(payload EnquiryMailPayloadUsingSendgrid, elapsed time.Duration, err error) error
func startLogConsumer(conn *amqp.Connection)
func startMailConsumer(conn *amqp.Connection)
func startEnquiryMailConsumer(conn *amqp.Connection)
```
> ### Logger-service
##### HTTP
```go
func (apiConfig *LocalApiConfig) WriteLog(c *gin.Context)
```
##### gRPC
```go
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error)
func (l *LogServer) GetAllLogs(ctx context.Context, request *logs.GetAllLogsRequest) (*logs.GetAllLogsResponse, error)
func GRPCListener(localApiConfig *LocalApiConfig)
```
> ### Mail-service
##### HTTP
```go
func (app *Config) SendMailViaSendGrid(w http.ResponseWriter, r *http.Request)
```

### Diagram
![img.png](img.png)


























