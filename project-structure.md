```
go_backend_gin/
│
├── main.go                     # Entry point of the application
│
│
├── controllers/                # HTTP request handlers
│   ├── patient_controller.go   # Patient-related request handling
│   └── staff_controller.go     # Staff-related request handling
│
├── database/
│   ├── connectDB.go     # connect database
│   
├── services/                   # Business logic layer
│   ├── patient_service.go      # Patient service logic
│   └── staff_service.go        # Staff service logic
│
├── repository/                 # Data access layer
│   ├── patient_repository.go   # Patient data access methods
│   └── staff_repository.go     # Staff data access methods
│
├── models/                     # Data models
│   ├── patient.go              # Patient model
│   └── staff.go                # Staff model
│
├── middleware/                 # Middleware functions
│   └── validateToken.go        # JWT token validation middleware
│
├── routes/                     # API route definitions
│   └── setuproutes.go          # Route setup and grouping
│
│
├── api_spec.md                 # API specification
├── database_schema.md          # Database schema documentation
│
├── tests/                      # Unit and integration tests
│
│
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency checksums
├── .env                        # Environment configuration
├── .gitignore                  # Git ignore file
└── README.md                   # Project README
```

### คำอธิบายโครงสร้างโปรเจ็ค

1. **main.go**: จุดเริ่มต้นของแอปพลิเคชัน

2. **controllers/**: 
   - จัดการ HTTP requests
   - แปลงข้อมูลจาก request เป็นรูปแบบที่เหมาะสม

3. **services/**: 
   - ประมวลผลตรรกะทางธุรกิจ
   - เชื่อมต่อระหว่าง controllers และ repositories

4. **repository/**: 
   - จัดการการเข้าถึงฐานข้อมูล
   - ดำเนินการ CRUD กับฐานข้อมูล

5. **models/**: 
   - นิยามโครงสร้างข้อมูล
   - ใช้กับ ORM (GORM)

6. **middleware/**: 
   - ฟังก์ชันสำหรับตรวจสอบ/กรองข้อมูล
   - ตรวจสอบ authentication

7. **routes/**: 
   - กำหนด API endpoints
   - จัดกลุ่ม routes

8. **tests/**: 
    - unit test


