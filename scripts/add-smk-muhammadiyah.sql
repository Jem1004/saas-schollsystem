-- Script untuk menambahkan SMK Muhammadiyah dan admin-nya
-- Jalankan script ini di database PostgreSQL
-- 
-- CARA MENJALANKAN:
-- 
-- Option 1: Menggunakan docker exec
-- docker exec -i school-management-postgres psql -U school_admin -d school_management < scripts/add-smk-muhammadiyah.sql
--
-- Option 2: Menggunakan psql langsung (jika PostgreSQL client terinstall)
-- psql -h localhost -U school_admin -d school_management -f scripts/add-smk-muhammadiyah.sql
-- Password: school_secret_2024
--
-- Option 3: Copy-paste ke Adminer (http://localhost:8080)
-- Server: postgres, Username: school_admin, Password: school_secret_2024, Database: school_management
--
-- Setelah menjalankan script ini, login dengan:
-- Username: admin.smkmuhammadiyah
-- Password: admin123

-- 1. Tambahkan sekolah SMK Muhammadiyah
INSERT INTO schools (name, address, phone, email, is_active, created_at, updated_at)
VALUES (
    'SMK Muhammadiyah',
    'Jl. Muhammadiyah No. 1',
    '021-9876543',
    'smkmuhammadiyah@edu.id',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (name) DO UPDATE SET updated_at = NOW()
RETURNING id;

-- 2. Tambahkan admin untuk SMK Muhammadiyah
-- Password: admin123 (bcrypt hash dari golang bcrypt.DefaultCost)
INSERT INTO users (school_id, role, username, password_hash, email, name, is_active, must_reset_pwd, created_at, updated_at)
SELECT 
    s.id,
    'admin_sekolah',
    'admin.smkmuhammadiyah',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.Q7Z5xQhqKa5lXUnMxK', -- admin123
    'admin@smkmuhammadiyah.edu.id',
    'Admin SMK Muhammadiyah',
    true,
    false,
    NOW(),
    NOW()
FROM schools s
WHERE s.name = 'SMK Muhammadiyah'
ON CONFLICT (username) DO UPDATE SET 
    school_id = EXCLUDED.school_id,
    updated_at = NOW();

-- 3. Tambahkan settings untuk SMK Muhammadiyah
INSERT INTO school_settings (school_id, attendance_start_time, attendance_end_time, attendance_late_threshold, attendance_very_late_threshold, enable_attendance_notification, enable_grade_notification, enable_bk_notification, enable_homeroom_notification, academic_year, semester, created_at, updated_at)
SELECT 
    s.id,
    '07:00',
    '07:30',
    30,
    60,
    true,
    true,
    true,
    true,
    '2024/2025',
    1,
    NOW(),
    NOW()
FROM schools s
WHERE s.name = 'SMK Muhammadiyah'
ON CONFLICT (school_id) DO UPDATE SET updated_at = NOW();

-- 4. Tambahkan beberapa kelas untuk SMK Muhammadiyah
INSERT INTO classes (school_id, name, grade, year, created_at, updated_at)
SELECT s.id, '10-TKJ-1', 10, '2024/2025', NOW(), NOW() FROM schools s WHERE s.name = 'SMK Muhammadiyah'
ON CONFLICT DO NOTHING;

INSERT INTO classes (school_id, name, grade, year, created_at, updated_at)
SELECT s.id, '10-TKJ-2', 10, '2024/2025', NOW(), NOW() FROM schools s WHERE s.name = 'SMK Muhammadiyah'
ON CONFLICT DO NOTHING;

INSERT INTO classes (school_id, name, grade, year, created_at, updated_at)
SELECT s.id, '11-TKJ-1', 11, '2024/2025', NOW(), NOW() FROM schools s WHERE s.name = 'SMK Muhammadiyah'
ON CONFLICT DO NOTHING;

-- Verifikasi data
SELECT '=== SCHOOLS ===' as info;
SELECT id, name, is_active FROM schools;

SELECT '=== ADMIN SEKOLAH USERS ===' as info;
SELECT u.id, u.username, u.role, u.school_id, s.name as school_name
FROM users u
LEFT JOIN schools s ON u.school_id = s.id
WHERE u.role = 'admin_sekolah';

SELECT '=== CLASSES ===' as info;
SELECT c.id, c.name, c.grade, s.name as school_name
FROM classes c
JOIN schools s ON c.school_id = s.id
ORDER BY s.name, c.grade, c.name;
