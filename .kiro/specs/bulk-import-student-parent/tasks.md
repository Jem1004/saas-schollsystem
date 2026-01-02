# Implementation Plan: Bulk Import Student & Parent

## Overview

Implementasi fitur bulk import untuk data siswa dan orang tua melalui file Excel. Menggunakan pendekatan hybrid untuk relasi kelas (auto-match jika cocok, warning jika tidak) dan import terpisah untuk orang tua dengan linking manual via UI.

## Tasks

- [x] 1. Modify Student Model for Nullable ClassID
  - [x] 1.1 Update Student model to make ClassID nullable
    - Change `ClassID uint` to `ClassID *uint` in student.go
    - Update Class relation to pointer `*Class`
    - Update Validate() to remove ClassID requirement
    - Add CanBeActive() method
    - _Requirements: 8.1, 8.2, 8.3_
  - [x] 1.2 Create database migration for ClassID nullable
    - Create migration script to alter students table
    - Update existing students to maintain data integrity
    - _Requirements: 8.1_
  - [x] 1.3 Update school service to handle nullable ClassID
    - Modify CreateStudent to accept optional ClassID
    - Modify UpdateStudent to handle ClassID changes
    - Add logic to set IsActive based on ClassID presence
    - _Requirements: 3.9, 3.10, 8.2, 8.3_
  - [ ]* 1.4 Write property test for ClassID-IsActive invariant
    - **Property 4: ClassID-IsActive Invariant**
    - **Validates: Requirements 3.9, 3.10, 6.4, 8.2, 8.3**

- [x] 2. Create Import Module Structure
  - [x] 2.1 Create import module directory and files
    - Create backend/internal/modules/import/
    - Create handler.go, service.go, dto.go, parser.go, matcher.go
    - _Requirements: 1.1, 1.2, 2.1_
  - [x] 2.2 Implement Excel parser using excelize library
    - Add excelize dependency to go.mod
    - Implement ParseStudentFile() and ParseParentFile()
    - Handle header row and data rows
    - _Requirements: 2.3_
  - [ ]* 2.3 Write unit tests for Excel parser
    - Test parsing valid files
    - Test handling empty files
    - Test handling invalid formats
    - _Requirements: 2.1, 2.2, 2.3, 2.4_

- [x] 3. Implement Template Generation
  - [x] 3.1 Implement GenerateStudentTemplate()
    - Create Excel file with columns: NIS, NISN, Nama, Kelas
    - Add header row and example data row
    - _Requirements: 1.1, 1.3, 1.4_
  - [x] 3.2 Implement GenerateParentTemplate()
    - Create Excel file with columns: Nama, No_HP, Email
    - Add header row and example data row
    - _Requirements: 1.2, 1.3, 1.4_
  - [x] 3.3 Create handler endpoints for template download
    - GET /api/v1/school/import/template/students
    - GET /api/v1/school/import/template/parents
    - _Requirements: 1.1, 1.2_

- [x] 4. Implement Student Import Service
  - [x] 4.1 Implement student row validation
    - Validate NIS, NISN, Name are not empty
    - Return appropriate error messages
    - _Requirements: 3.1, 3.2, 3.3_
  - [ ]* 4.2 Write property test for required field validation
    - **Property 1: Required Field Validation**
    - **Validates: Requirements 3.1, 3.2, 3.3**
  - [x] 4.3 Implement duplicate detection for NISN and NIS
    - Check existing students by NISN
    - Check existing students by NIS within school
    - Add to warning list if duplicate found
    - _Requirements: 3.4, 3.5_
  - [ ]* 4.4 Write property test for duplicate detection
    - **Property 2: Duplicate Detection**
    - **Validates: Requirements 3.4, 3.5**
  - [x] 4.5 Implement class matcher
    - Find class by name (case-insensitive)
    - Return class ID if found, nil if not
    - _Requirements: 3.6, 3.7_
  - [ ]* 4.6 Write property test for class matching
    - **Property 3: Class Matching Correctness**
    - **Validates: Requirements 3.6**
  - [x] 4.7 Implement ImportStudents() service method
    - Parse Excel file
    - Validate each row
    - Check duplicates
    - Match classes
    - Create student records with transaction
    - Set IsActive based on ClassID
    - Return ImportResult
    - _Requirements: 3.1-3.11, 5.1-5.6_
  - [ ]* 4.8 Write property test for tenant isolation
    - **Property 5: Tenant Isolation**
    - **Validates: Requirements 3.11, 4.6, 6.5, 7.6**
  - [ ]* 4.9 Write property test for import result accuracy
    - **Property 7: Import Result Accuracy**
    - **Validates: Requirements 5.1, 5.2, 5.3, 5.4**

- [x] 5. Implement Parent Import Service
  - [x] 5.1 Implement parent row validation
    - Validate Name is not empty
    - Validate Email format
    - _Requirements: 4.1, 4.2_
  - [x] 5.2 Implement duplicate detection for parent email/phone
    - Check existing users by phone (username)
    - Add to warning list if duplicate found
    - _Requirements: 4.3_
  - [x] 5.3 Implement ImportParents() service method
    - Parse Excel file
    - Validate each row
    - Check duplicates
    - Create User account with role orang_tua
    - Create Parent record
    - Generate temporary password
    - Return ImportResult
    - _Requirements: 4.1-4.7, 5.1-5.6_
  - [ ]* 5.4 Write property test for parent account creation
    - **Property 6: Parent Account Creation Consistency**
    - **Validates: Requirements 4.4, 4.5**

- [x] 6. Create Import Handler and Routes
  - [x] 6.1 Implement import handler
    - POST /api/v1/school/import/students
    - POST /api/v1/school/import/parents
    - Handle file upload (multipart/form-data)
    - Validate file size (max 5MB)
    - Return ImportResult
    - _Requirements: 2.1, 2.2, 2.5_
  - [x] 6.2 Register import routes in main router
    - Add import module to dependency injection
    - Register routes with auth middleware
    - _Requirements: 2.1_

- [x] 7. Checkpoint - Backend Import Complete
  - Ensure all backend tests pass
  - Test import endpoints manually with sample Excel files
  - Ask the user if questions arise

- [x] 8. Implement Bulk Class Assignment
  - [x] 8.1 Add GetStudentsWithoutClass() to school service
    - Query students where ClassID is null
    - Return list of students needing class assignment
    - _Requirements: 6.1_
  - [x] 8.2 Add BulkAssignClass() to school service
    - Accept list of student IDs and class ID
    - Validate class belongs to same school
    - Update ClassID for all students
    - Set IsActive to true
    - _Requirements: 6.2, 6.3, 6.4, 6.5_
  - [ ]* 8.3 Write property test for bulk class assignment
    - **Property 9: Bulk Class Assignment Completeness**
    - **Validates: Requirements 6.3, 6.4**
  - [x] 8.4 Create handler endpoint for bulk class assignment
    - POST /api/v1/school/students/bulk-assign-class
    - _Requirements: 6.2_

- [x] 9. Enhance Parent-Student Linking
  - [x] 9.1 Add search students endpoint for linking
    - GET /api/v1/school/students/search?q={query}
    - Search by NISN or name
    - _Requirements: 7.2_
  - [ ]* 9.2 Write property test for parent-student linking
    - **Property 8: Parent-Student Many-to-Many Relationship**
    - **Validates: Requirements 7.3, 7.4, 7.5**

- [x] 10. Checkpoint - Backend Complete
  - Ensure all backend tests pass
  - Verify all API endpoints work correctly
  - Ask the user if questions arise

- [x] 11. Frontend - Student Import UI
  - [x] 11.1 Add import service functions
    - Create web-admin/src/services/import.ts
    - Add downloadStudentTemplate()
    - Add importStudents()
    - _Requirements: 1.1, 3.1-3.11_
  - [x] 11.2 Update StudentManagement.vue with import UI
    - Add "Download Template" button
    - Add "Import Excel" button with file upload
    - Show import result dialog (success, errors, warnings)
    - Show count of students needing class assignment
    - _Requirements: 1.1, 5.1-5.5_
  - [x] 11.3 Add bulk class assignment UI
    - Add filter for students without class
    - Add checkbox selection for multiple students
    - Add "Assign Class" button with class selector
    - _Requirements: 6.1, 6.2_

- [x] 12. Frontend - Parent Import UI
  - [x] 12.1 Add parent import service functions
    - Add downloadParentTemplate()
    - Add importParents()
    - _Requirements: 1.2, 4.1-4.7_
  - [x] 12.2 Update ParentManagement.vue with import UI
    - Add "Download Template" button
    - Add "Import Excel" button with file upload
    - Show import result dialog
    - _Requirements: 1.2, 5.1-5.5_
  - [x] 12.3 Enhance parent-student linking UI
    - Add student search in parent detail/edit dialog
    - Allow search by NISN or name
    - Show linked students list
    - _Requirements: 7.1, 7.2, 7.3_

- [x] 13. Final Checkpoint
  - Ensure all tests pass
  - Test complete import workflow end-to-end
  - Verify UI displays correct feedback
  - Ask the user if questions arise

## Notes

- Tasks marked with `*` are optional property-based tests
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties
- Unit tests validate specific examples and edge cases
- Use `excelize` library for Excel file handling in Go
- Use `gopter` library for property-based testing in Go
