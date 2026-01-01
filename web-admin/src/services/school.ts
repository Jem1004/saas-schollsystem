import api from './api'
import type {
  Class,
  Student,
  Parent,
  SchoolUser,
  SchoolSettings,
  SchoolStats,
  AttendanceSummary,
  CreateClassRequest,
  UpdateClassRequest,
  CreateStudentRequest,
  UpdateStudentRequest,
  CreateParentRequest,
  UpdateParentRequest,
  CreateUserRequest,
  UpdateUserRequest,
  UpdateSchoolSettingsRequest,
  ClassListResponse,
  StudentListResponse,
  ParentListResponse,
  UserListResponse,
  AttendanceListResponse,
  ResetPasswordResponse,
  SchoolAttendanceSummaryResponse,
  Device,
  PairingSessionResponse,
} from '@/types/school'

export interface PaginationParams {
  page?: number
  page_size?: number
  search?: string
  class_id?: number
  date?: string
  start_date?: string
  end_date?: string
}

// API response wrapper
interface ApiResponse<T> {
  success: boolean
  data: T
  message?: string
}

// Backend response types (snake_case)
interface StudentApiResponse {
  id: number
  school_id: number
  class_id: number
  class_name?: string
  nis: string
  nisn: string
  name: string
  rfid_code?: string
  is_active: boolean
  has_account?: boolean
  username?: string
  temporary_password?: string
  created_at: string
  updated_at: string
}

interface StudentListApiResponse {
  students: StudentApiResponse[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

// Transform student from API response to frontend format
function transformStudent(apiStudent: StudentApiResponse): Student {
  return {
    id: apiStudent.id,
    schoolId: apiStudent.school_id,
    classId: apiStudent.class_id,
    className: apiStudent.class_name,
    nis: apiStudent.nis,
    nisn: apiStudent.nisn,
    name: apiStudent.name,
    rfidCode: apiStudent.rfid_code,
    isActive: apiStudent.is_active,
    hasAccount: apiStudent.has_account,
    username: apiStudent.username,
    temporaryPassword: apiStudent.temporary_password,
    createdAt: apiStudent.created_at,
    updatedAt: apiStudent.updated_at,
  }
}

// Transform student list response
function transformStudentList(apiResponse: StudentListApiResponse): StudentListResponse {
  return {
    students: apiResponse.students.map(transformStudent),
    pagination: {
      page: apiResponse.pagination.page,
      page_size: apiResponse.pagination.page_size,
      total: apiResponse.pagination.total,
      total_pages: apiResponse.pagination.total_pages,
    },
  }
}

// Backend response types for Class (snake_case)
interface ClassApiResponse {
  id: number
  school_id: number
  name: string
  grade: number
  year: string
  homeroom_teacher_id?: number
  homeroom_teacher?: {
    id: number
    name: string
    username: string
  }
  student_count: number
  created_at: string
  updated_at: string
}

interface ClassListApiResponse {
  classes: ClassApiResponse[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

// Transform class from API response to frontend format
function transformClass(apiClass: ClassApiResponse): Class {
  return {
    id: apiClass.id,
    schoolId: apiClass.school_id,
    name: apiClass.name,
    grade: apiClass.grade,
    year: apiClass.year,
    homeroomTeacherId: apiClass.homeroom_teacher_id,
    homeroomTeacherName: apiClass.homeroom_teacher?.name,
    studentCount: apiClass.student_count,
    createdAt: apiClass.created_at,
    updatedAt: apiClass.updated_at,
  }
}

// Transform class list response
function transformClassList(apiResponse: ClassListApiResponse): ClassListResponse {
  return {
    classes: apiResponse.classes.map(transformClass),
    pagination: {
      page: apiResponse.pagination.page,
      page_size: apiResponse.pagination.page_size,
      total: apiResponse.pagination.total,
      total_pages: apiResponse.pagination.total_pages,
    },
  }
}

// Backend response types for Parent (snake_case)
interface ParentApiResponse {
  id: number
  school_id: number
  user_id: number
  name: string
  phone?: string
  email?: string
  username?: string
  temporary_password?: string
  students?: StudentApiResponse[]
  created_at: string
  updated_at: string
}

interface ParentListApiResponse {
  parents: ParentApiResponse[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

// Transform parent from API response to frontend format
function transformParent(apiParent: ParentApiResponse): Parent {
  return {
    id: apiParent.id,
    schoolId: apiParent.school_id,
    userId: apiParent.user_id,
    name: apiParent.name,
    phone: apiParent.phone,
    email: apiParent.email,
    username: apiParent.username,
    temporaryPassword: apiParent.temporary_password,
    studentIds: apiParent.students?.map(s => s.id) || [],
    studentNames: apiParent.students?.map(s => s.name) || [],
    createdAt: apiParent.created_at,
    updatedAt: apiParent.updated_at,
  }
}

// Transform parent list response
function transformParentList(apiResponse: ParentListApiResponse): ParentListResponse {
  return {
    parents: apiResponse.parents.map(transformParent),
    pagination: {
      page: apiResponse.pagination.page,
      page_size: apiResponse.pagination.page_size,
      total: apiResponse.pagination.total,
      total_pages: apiResponse.pagination.total_pages,
    },
  }
}

// Backend response type for ResetPassword (snake_case)
interface ResetPasswordApiResponse {
  username: string
  temporary_password: string
  message: string
}

// Transform reset password response
function transformResetPassword(apiResponse: ResetPasswordApiResponse): ResetPasswordResponse {
  return {
    username: apiResponse.username,
    temporaryPassword: apiResponse.temporary_password,
    message: apiResponse.message,
  }
}

// Backend response types for Device (snake_case)
interface DeviceApiResponse {
  id: number
  school_id: number
  device_code: string
  description: string
  is_active: boolean
  last_seen_at?: string
  created_at: string
  updated_at: string
}

// Transform device from API response to frontend format
function transformDevice(apiDevice: DeviceApiResponse): Device {
  return {
    id: apiDevice.id,
    schoolId: apiDevice.school_id,
    deviceCode: apiDevice.device_code,
    description: apiDevice.description,
    isActive: apiDevice.is_active,
    lastSeenAt: apiDevice.last_seen_at,
    createdAt: apiDevice.created_at,
    updatedAt: apiDevice.updated_at,
  }
}

// Backend response types for Pairing (snake_case)
interface PairingSessionApiResponse {
  active: boolean
  student_id?: number
  student_name?: string
  device_id?: number
  expires_at?: string
  message: string
}

// Transform pairing session from API response to frontend format
function transformPairingSession(apiResponse: PairingSessionApiResponse): PairingSessionResponse {
  return {
    active: apiResponse.active,
    studentId: apiResponse.student_id,
    studentName: apiResponse.student_name,
    deviceId: apiResponse.device_id,
    expiresAt: apiResponse.expires_at,
    message: apiResponse.message,
  }
}

// Backend response types for User (snake_case)
interface AssignedClassApiInfo {
  id: number
  name: string
}

interface UserApiResponse {
  id: number
  school_id: number
  role: 'guru' | 'wali_kelas' | 'guru_bk'
  username: string
  email?: string
  name?: string
  is_active: boolean
  must_reset_pwd: boolean
  assigned_class_id?: number      // For wali_kelas
  assigned_class_name?: string    // For wali_kelas
  assigned_classes?: AssignedClassApiInfo[]  // For guru_bk
  last_login_at?: string
  created_at: string
  updated_at: string
}

interface UserListApiResponse {
  users: UserApiResponse[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

// Transform user from API response to frontend format
function transformUser(apiUser: UserApiResponse): SchoolUser {
  return {
    id: apiUser.id,
    schoolId: apiUser.school_id,
    role: apiUser.role,
    username: apiUser.username,
    email: apiUser.email,
    name: apiUser.name,
    isActive: apiUser.is_active,
    mustResetPwd: apiUser.must_reset_pwd,
    assignedClassId: apiUser.assigned_class_id,
    assignedClassName: apiUser.assigned_class_name,
    assignedClasses: apiUser.assigned_classes,
    lastLoginAt: apiUser.last_login_at,
    createdAt: apiUser.created_at,
    updatedAt: apiUser.updated_at,
  }
}

// Transform user list response
function transformUserList(apiResponse: UserListApiResponse): UserListResponse {
  return {
    users: apiResponse.users.map(transformUser),
    pagination: {
      page: apiResponse.pagination.page,
      page_size: apiResponse.pagination.page_size,
      total: apiResponse.pagination.total,
      total_pages: apiResponse.pagination.total_pages,
    },
  }
}

export const schoolService = {
  // Stats
  async getStats(): Promise<SchoolStats> {
    const response = await api.get<ApiResponse<SchoolStats>>('/school/stats')
    return response.data.data
  },

  // Classes
  async getClasses(params?: PaginationParams): Promise<ClassListResponse> {
    const response = await api.get<ApiResponse<ClassListApiResponse>>('/school/classes', { params })
    return transformClassList(response.data.data)
  },

  async getClass(id: number): Promise<Class> {
    const response = await api.get<ApiResponse<ClassApiResponse>>(`/school/classes/${id}`)
    return transformClass(response.data.data)
  },

  async createClass(data: CreateClassRequest): Promise<Class> {
    const response = await api.post<ApiResponse<ClassApiResponse>>('/school/classes', data)
    return transformClass(response.data.data)
  },

  async updateClass(id: number, data: UpdateClassRequest): Promise<Class> {
    const response = await api.put<ApiResponse<ClassApiResponse>>(`/school/classes/${id}`, data)
    return transformClass(response.data.data)
  },

  async deleteClass(id: number): Promise<void> {
    await api.delete(`/school/classes/${id}`)
  },

  // Students
  async getStudents(params?: PaginationParams): Promise<StudentListResponse> {
    const response = await api.get<ApiResponse<StudentListApiResponse>>('/school/students', { params })
    return transformStudentList(response.data.data)
  },

  async getStudent(id: number): Promise<Student> {
    const response = await api.get<ApiResponse<StudentApiResponse>>(`/school/students/${id}`)
    return transformStudent(response.data.data)
  },

  async createStudent(data: CreateStudentRequest): Promise<Student> {
    const response = await api.post<ApiResponse<StudentApiResponse>>('/school/students', data)
    return transformStudent(response.data.data)
  },

  async updateStudent(id: number, data: UpdateStudentRequest): Promise<Student> {
    const response = await api.put<ApiResponse<StudentApiResponse>>(`/school/students/${id}`, data)
    return transformStudent(response.data.data)
  },

  async deleteStudent(id: number): Promise<void> {
    await api.delete(`/school/students/${id}`)
  },

  async createStudentAccount(id: number): Promise<Student> {
    const response = await api.post<ApiResponse<StudentApiResponse>>(`/school/students/${id}/account`)
    return transformStudent(response.data.data)
  },

  async resetStudentPassword(id: number): Promise<ResetPasswordResponse> {
    const response = await api.post<ApiResponse<ResetPasswordApiResponse>>(`/school/students/${id}/reset-password`)
    return transformResetPassword(response.data.data)
  },

  // Parents
  async getParents(params?: PaginationParams): Promise<ParentListResponse> {
    const response = await api.get<ApiResponse<ParentListApiResponse>>('/school/parents', { params })
    return transformParentList(response.data.data)
  },

  async getParent(id: number): Promise<Parent> {
    const response = await api.get<ApiResponse<ParentApiResponse>>(`/school/parents/${id}`)
    return transformParent(response.data.data)
  },

  async createParent(data: CreateParentRequest): Promise<Parent> {
    const response = await api.post<ApiResponse<ParentApiResponse>>('/school/parents', data)
    return transformParent(response.data.data)
  },

  async updateParent(id: number, data: UpdateParentRequest): Promise<Parent> {
    const response = await api.put<ApiResponse<ParentApiResponse>>(`/school/parents/${id}`, data)
    return transformParent(response.data.data)
  },

  async deleteParent(id: number): Promise<void> {
    await api.delete(`/school/parents/${id}`)
  },

  async linkParentToStudent(parentId: number, studentId: number): Promise<void> {
    await api.post(`/school/parents/${parentId}/students/${studentId}`)
  },

  async unlinkParentFromStudent(parentId: number, studentId: number): Promise<void> {
    await api.delete(`/school/parents/${parentId}/students/${studentId}`)
  },

  async resetParentPassword(id: number): Promise<ResetPasswordResponse> {
    const response = await api.post<ApiResponse<ResetPasswordApiResponse>>(`/school/parents/${id}/reset-password`)
    return transformResetPassword(response.data.data)
  },

  // Users
  async getUsers(params?: PaginationParams): Promise<UserListResponse> {
    const response = await api.get<ApiResponse<UserListApiResponse>>('/school/users', { params })
    return transformUserList(response.data.data)
  },

  async getUser(id: number): Promise<SchoolUser> {
    const response = await api.get<ApiResponse<UserApiResponse>>(`/school/users/${id}`)
    return transformUser(response.data.data)
  },

  async createUser(data: CreateUserRequest): Promise<SchoolUser> {
    const response = await api.post<ApiResponse<UserApiResponse>>('/school/users', data)
    return transformUser(response.data.data)
  },

  async updateUser(id: number, data: UpdateUserRequest): Promise<SchoolUser> {
    const response = await api.put<ApiResponse<UserApiResponse>>(`/school/users/${id}`, data)
    return transformUser(response.data.data)
  },

  async deleteUser(id: number): Promise<void> {
    await api.delete(`/school/users/${id}`)
  },

  async resetUserPassword(id: number): Promise<ResetPasswordResponse> {
    const response = await api.post<ApiResponse<ResetPasswordApiResponse>>(`/school/users/${id}/reset-password`)
    return transformResetPassword(response.data.data)
  },

  // Attendance
  async getAttendanceSummary(params?: PaginationParams): Promise<AttendanceListResponse> {
    const response = await api.get<ApiResponse<SchoolAttendanceSummaryResponse>>('/attendance/summary', { params })
    const data = response.data.data
    
    // Transform backend response to frontend format
    const attendanceData: AttendanceSummary[] = data.by_class?.map(item => ({
      date: data.date,
      classId: item.class_id,
      className: item.class_name,
      totalStudents: item.total_students,
      present: item.present,
      late: item.late,
      absent: item.absent,
    })) || []

    return {
      data: attendanceData,
      total: attendanceData.length,
      page: 1,
      pageSize: attendanceData.length,
    }
  },

  async getAttendanceByClass(classId: number, date: string): Promise<AttendanceSummary> {
    const response = await api.get<ApiResponse<AttendanceSummary>>(`/attendance/class/${classId}`, { params: { date } })
    return response.data.data
  },

  async exportAttendance(params: { start_date: string; end_date: string; class_id?: number }): Promise<Blob> {
    const response = await api.get('/attendance/export', {
      params,
      responseType: 'blob',
    })
    return response.data
  },

  // Settings
  async getSettings(): Promise<SchoolSettings> {
    const response = await api.get<ApiResponse<SchoolSettings>>('/settings')
    return response.data.data
  },

  async updateSettings(data: UpdateSchoolSettingsRequest): Promise<SchoolSettings> {
    const response = await api.put<ApiResponse<SchoolSettings>>('/settings', data)
    return response.data.data
  },

  async resetSettings(): Promise<SchoolSettings> {
    const response = await api.post<ApiResponse<SchoolSettings>>('/settings/reset')
    return response.data.data
  },

  // Devices
  async getSchoolDevices(): Promise<Device[]> {
    const response = await api.get<ApiResponse<DeviceApiResponse[]>>('/school/devices')
    return response.data.data.map(transformDevice)
  },

  // RFID Pairing
  async startPairing(deviceId: number, studentId: number): Promise<PairingSessionResponse> {
    const response = await api.post<ApiResponse<PairingSessionApiResponse>>('/pairing/start', {
      device_id: deviceId,
      student_id: studentId,
    })
    return transformPairingSession(response.data.data)
  },

  async cancelPairing(deviceId: number): Promise<void> {
    await api.post(`/pairing/cancel/${deviceId}`)
  },

  async getPairingStatus(deviceId: number): Promise<PairingSessionResponse> {
    const response = await api.get<ApiResponse<PairingSessionApiResponse>>(`/pairing/status/${deviceId}`)
    return transformPairingSession(response.data.data)
  },

  async clearStudentRFID(studentId: number): Promise<void> {
    await api.post(`/school/students/${studentId}/clear-rfid`)
  },
}
