# ข้อกำหนดของ API ระบบจัดการข้อมูลผู้ป่วยและบุคลากร

## 1. การรับรองความถูกต้อง (Authentication)

### 1.1 เข้าสู่ระบบสำหรับบุคลากร
- **เส้นทาง**: `/staff/login`
- **วิธี**: POST
- **คำขอ**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **การตอบกลับ**:
  - สำเร็จ: 
    ```json
    {
      "status": "Login success",
      "token": "jwt_token"
    }
    ```
  - ล้มเหลว: 
    ```json
    {
      "error": "ข้อความแสดงข้อผิดพลาด"
    }
    ```

### 1.2 สร้างบัญชีบุคลากร
- **เส้นทาง**: `/staff/create`
- **วิธี**: POST
- **ต้องการ Token**: ใช่
- **คำขอ**:
  ```json
  {
    "username": "string",
    "password": "string",
    "hospital_id": "string"
  }
  ```
- **การตอบกลับ**:
  - สำเร็จ: 
    ```json
    {
      "status": "Create Staff ID success"
    }
    ```
  - ล้มเหลว: 
    ```json
    {
      "error": "ข้อความแสดงข้อผิดพลาด"
    }
    ```

## 2. การค้นหาข้อมูลผู้ป่วย

### 2.1 ค้นหาผู้ป่วยแบบละเอียด
- **เส้นทาง**: `/patient/search`
- **วิธี**: POST
- **ต้องการ Token**: ใช่
- **คำขอ**:
  ```json
  {
    "national_id": "string (optional)",
    "passport_id": "string (optional)",
    "first_name": "string (optional)",
    "middle_name": "string (optional)",
    "last_name": "string (optional)",
    "date_of_birth": "string (optional)",
    "phone_number": "string (optional)",
    "email": "string (optional)",
    "hospital_id": "integer (required)"
  }
  ```
- **การตอบกลับ**:
  - สำเร็จ: 
    ```json
    {
      "patients": [
        {
          "patient_hn": "string",
          "first_name_th": "string",
          "middle_name_th": "string",
          "last_name_th": "string",
          "date_of_birth": "string",
          "national_id": "string",
          "passport_id": "string",
          "phone_number": "string",
          "email": "string",
          "gender": "string"
        }
      ]
    }
    ```
  - ล้มเหลว: 
    ```json
    {
      "error": "ข้อความแสดงข้อผิดพลาด"
    }
    ```

### 2.2 ค้นหาผู้ป่วยด้วย ID
- **เส้นทาง**: `/patient/search/:id`
- **วิธี**: GET
- **พารามิเตอร์**: `id` (national_id หรือ passport_id)
- **การตอบกลับ**:
  - สำเร็จ: 
    ```json
    {
      "patients": [
        {
          "patient_hn": "string",
          "first_name_th": "string",
          "middle_name_th": "string",
          "last_name_th": "string",
          "date_of_birth": "string",
          "national_id": "string",
          "passport_id": "string",
          "phone_number": "string",
          "email": "string",
          "gender": "string"
        }
      ]
    }
    ```
  - ล้มเหลว: 
    ```json
    {
      "error": "ข้อความแสดงข้อผิดพลาด"
    }
    ```

## ข้อควรระวัง
- ทุก API ที่ต้องการ Token จะต้องส่ง Header Authorization เป็น `Bearer {token}`
- Token มีอายุการใช้งาน 24 ชั่วโมง
- การค้นหาผู้ป่วยจำเป็นต้องระบุ hospital_id
- ต้องระบุอย่างน้อยหนึ่งเกณฑ์ในการค้นหาผู้ป่วย

## การจัดการข้อผิดพลาด
- รหัสสถานะ HTTP 200: สำหรับการตอบกลับปกติ
- รหัสสถานะ HTTP 400: ข้อมูลที่ส่งมาไม่ถูกต้อง
- รหัสสถานะ HTTP 401: ไม่ได้รับอนุญาต (token ไม่ถูกต้อง)
