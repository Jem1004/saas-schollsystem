# Requirements Document

## Introduction

Fitur Bulk Import untuk data Siswa dan Orang Tua memungkinkan Admin Sekolah mengimpor data secara massal melalui file Excel. Fitur ini mendukung workflow hybrid untuk relasi kelas (auto-match jika cocok, warning jika tidak) dan import terpisah untuk orang tua dengan linking manual via UI.

## Glossary

- **Import_Service**: Layanan backend yang memproses file Excel dan menyimpan data ke database
- **Excel_Parser**: Komponen yang membaca dan memvalidasi file Excel (.xlsx)
- **Student_Importer**: Komponen yang mengimpor data siswa dengan logic class matching
- **Parent_Importer**: Komponen yang mengimpor data orang tua tanpa linking
- **Import_Result**: Objek yang berisi hasil import (sukses, gagal, warning)
- **Class_Matcher**: Komponen yang mencocokkan nama kelas dari file dengan data kelas di database
- **Student_Parent_Linker**: Komponen UI untuk menghubungkan orang tua dengan siswa
- **NIS**: Nomor Induk Siswa (unik per sekolah)
- **NISN**: Nomor Induk Siswa Nasional (unik nasional)

## Requirements

### Requirement 1: Download Template Excel

**User Story:** As an Admin Sekolah, I want to download Excel templates for student and parent import, so that I can prepare data in the correct format.

#### Acceptance Criteria

1. WHEN an admin requests student template, THE Import_Service SHALL generate an Excel file with columns: NIS, NISN, Nama, Kelas
2. WHEN an admin requests parent template, THE Import_Service SHALL generate an Excel file with columns: Nama, No_HP, Email
3. THE Excel_Parser SHALL include header row with column names in the template
4. THE Excel_Parser SHALL include example data row in the template to guide users

### Requirement 2: Upload and Parse Excel File

**User Story:** As an Admin Sekolah, I want to upload an Excel file for import, so that the system can process my data.

#### Acceptance Criteria

1. WHEN an admin uploads a file, THE Import_Service SHALL validate the file format is .xlsx
2. IF the file format is invalid, THEN THE Import_Service SHALL return an error message specifying accepted formats
3. WHEN parsing the Excel file, THE Excel_Parser SHALL read all rows starting from row 2 (after header)
4. IF the Excel file is empty or has no data rows, THEN THE Import_Service SHALL return an error indicating no data found
5. WHEN the file size exceeds 5MB, THE Import_Service SHALL reject the upload with appropriate error message

### Requirement 3: Import Student Data

**User Story:** As an Admin Sekolah, I want to import student data from Excel, so that I can register multiple students efficiently.

#### Acceptance Criteria

1. WHEN importing a student row, THE Student_Importer SHALL validate that NIS is not empty
2. WHEN importing a student row, THE Student_Importer SHALL validate that NISN is not empty
3. WHEN importing a student row, THE Student_Importer SHALL validate that Nama is not empty
4. IF a student with the same NISN already exists in the school, THEN THE Student_Importer SHALL skip the row and add to warning list
5. IF a student with the same NIS already exists in the school, THEN THE Student_Importer SHALL skip the row and add to warning list
6. WHEN class_name is provided and matches an existing class, THE Class_Matcher SHALL assign the ClassID to the student
7. WHEN class_name is provided but does not match any existing class, THE Class_Matcher SHALL set ClassID to null and add warning to result
8. WHEN class_name is empty, THE Student_Importer SHALL set ClassID to null
9. WHEN a student is imported without ClassID, THE Student_Importer SHALL set IsActive to false
10. WHEN a student is imported with valid ClassID, THE Student_Importer SHALL set IsActive to true
11. THE Student_Importer SHALL associate all imported students with the admin's school (tenant isolation)

### Requirement 4: Import Parent Data

**User Story:** As an Admin Sekolah, I want to import parent data from Excel, so that I can register multiple parents efficiently.

#### Acceptance Criteria

1. WHEN importing a parent row, THE Parent_Importer SHALL validate that Nama is not empty
2. WHEN importing a parent row, THE Parent_Importer SHALL validate that Email is not empty and valid format
3. IF a parent with the same email already exists in the school, THEN THE Parent_Importer SHALL skip the row and add to warning list
4. WHEN a parent is imported, THE Parent_Importer SHALL create a User account with role "orang_tua"
5. WHEN creating parent user account, THE Import_Service SHALL generate a temporary password
6. THE Parent_Importer SHALL associate all imported parents with the admin's school (tenant isolation)
7. THE Parent_Importer SHALL NOT create any student links during import

### Requirement 5: Import Result and Feedback

**User Story:** As an Admin Sekolah, I want to see detailed import results, so that I can understand what was imported and what needs attention.

#### Acceptance Criteria

1. WHEN import completes, THE Import_Service SHALL return total rows processed count
2. WHEN import completes, THE Import_Service SHALL return successful imports count
3. WHEN import completes, THE Import_Service SHALL return failed imports count with row numbers and error messages
4. WHEN import completes, THE Import_Service SHALL return warning list with row numbers and warning messages
5. WHEN student import has students without class assignment, THE Import_Result SHALL include count of students needing class assignment
6. THE Import_Service SHALL process import within database transaction (all or nothing for valid rows)

### Requirement 6: Bulk Class Assignment

**User Story:** As an Admin Sekolah, I want to assign classes to multiple students at once, so that I can efficiently manage students imported without class.

#### Acceptance Criteria

1. WHEN admin views student list, THE System SHALL provide filter for students without class assignment
2. WHEN admin selects multiple students, THE System SHALL allow bulk class assignment
3. WHEN bulk class assignment is performed, THE System SHALL update ClassID for all selected students
4. WHEN a student receives class assignment, THE System SHALL set IsActive to true
5. THE System SHALL validate that selected class belongs to the same school

### Requirement 7: Link Parent to Student

**User Story:** As an Admin Sekolah, I want to link parents to their children, so that parents can access their children's information.

#### Acceptance Criteria

1. WHEN admin views parent detail, THE System SHALL show option to link students
2. WHEN linking student to parent, THE System SHALL allow search by NISN or student name
3. WHEN a link is created, THE Student_Parent_Linker SHALL create entry in student_parents table
4. THE System SHALL allow one parent to be linked to multiple students
5. THE System SHALL allow one student to be linked to multiple parents
6. THE System SHALL validate that linked student belongs to the same school

### Requirement 8: Model Changes for Nullable ClassID

**User Story:** As a system, I need to support students without class assignment, so that import workflow can function properly.

#### Acceptance Criteria

1. THE Student model SHALL allow ClassID to be null
2. WHEN ClassID is null, THE Student SHALL have IsActive set to false
3. WHEN ClassID is set from null to valid value, THE System SHALL allow IsActive to be set to true
4. THE System SHALL prevent attendance recording for students with null ClassID
