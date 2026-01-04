/// Application strings in Indonesian for the School Management mobile app.
///
/// All user-facing text is centralized here for easy maintenance and
/// potential future localization support.
///
/// Usage:
/// ```dart
/// Text(AppStrings.login)
/// Text(AppStrings.errorInvalidCredentials)
/// ```
class AppStrings {
  AppStrings._();

  // ============================================
  // App Info
  // ============================================

  static const String appName = 'School Management';
  static const String appTagline = 'Sistem Manajemen Sekolah';

  // ============================================
  // Common Actions
  // ============================================

  static const String ok = 'OK';
  static const String cancel = 'Batal';
  static const String save = 'Simpan';
  static const String delete = 'Hapus';
  static const String edit = 'Edit';
  static const String close = 'Tutup';
  static const String back = 'Kembali';
  static const String next = 'Selanjutnya';
  static const String done = 'Selesai';
  static const String retry = 'Coba Lagi';
  static const String refresh = 'Muat Ulang';
  static const String search = 'Cari';
  static const String filter = 'Filter';
  static const String seeAll = 'Lihat Semua';
  static const String seeMore = 'Lihat Lebih';
  static const String seeLess = 'Lihat Lebih Sedikit';
  static const String loading = 'Memuat...';
  static const String submit = 'Kirim';
  static const String confirm = 'Konfirmasi';
  static const String yes = 'Ya';
  static const String no = 'Tidak';

  // ============================================
  // Authentication
  // ============================================

  static const String login = 'Masuk';
  static const String logout = 'Keluar';
  static const String logoutConfirmTitle = 'Keluar dari Aplikasi';
  static const String logoutConfirmMessage = 'Apakah Anda yakin ingin keluar?';
  static const String welcomeBack = 'Selamat Datang';
  static const String loginSubtitle = 'Masuk untuk melanjutkan';
  static const String nisn = 'NISN';
  static const String nisnHint = 'Masukkan NISN';
  static const String password = 'Password';
  static const String passwordHint = 'Masukkan password';
  static const String forgotPassword = 'Lupa password?';
  static const String forgotPasswordHint = 'Hubungi admin sekolah';
  static const String changePassword = 'Ubah Password';
  static const String oldPassword = 'Password Lama';
  static const String newPassword = 'Password Baru';
  static const String confirmPassword = 'Konfirmasi Password';
  static const String passwordChanged = 'Password berhasil diubah';
  static const String mustChangePassword = 'Anda harus mengubah password untuk melanjutkan';

  // ============================================
  // Navigation
  // ============================================

  static const String home = 'Beranda';
  static const String dashboard = 'Dashboard';
  static const String attendance = 'Kehadiran';
  static const String grades = 'Nilai';
  static const String notes = 'Catatan';
  static const String bk = 'BK';
  static const String notifications = 'Notifikasi';
  static const String settings = 'Pengaturan';
  static const String profile = 'Profil';

  // ============================================
  // Greetings
  // ============================================

  static const String goodMorning = 'Selamat Pagi';
  static const String goodAfternoon = 'Selamat Siang';
  static const String goodEvening = 'Selamat Sore';
  static const String goodNight = 'Selamat Malam';

  // ============================================
  // Child Selector
  // ============================================

  static const String selectChild = 'Pilih Anak';
  static const String switchChild = 'Ganti Anak';
  static const String childClass = 'Kelas';

  // ============================================
  // Dashboard
  // ============================================

  static const String todayStatus = 'Status Hari Ini';
  static const String checkIn = 'Masuk';
  static const String checkOut = 'Pulang';
  static const String notYet = 'Belum';
  static const String summary = 'Ringkasan';
  static const String recentActivities = 'Aktivitas Terbaru';
  static const String attendancePercentage = 'Kehadiran';
  static const String averageGrade = 'Rata-rata Nilai';
  static const String achievementPoints = 'Poin Prestasi';
  static const String violationPoints = 'Poin Pelanggaran';
  static const String netPoints = 'Poin Bersih';

  // ============================================
  // Attendance
  // ============================================

  static const String attendanceHistory = 'Riwayat Kehadiran';
  static const String attendanceSummary = 'Ringkasan Kehadiran';
  static const String onTime = 'Tepat Waktu';
  static const String late = 'Terlambat';
  static const String veryLate = 'Sangat Terlambat';
  static const String absent = 'Tidak Hadir';
  static const String sick = 'Sakit';
  static const String excused = 'Izin';
  static const String present = 'Hadir';
  static const String totalDays = 'Total Hari';
  static const String filterByDate = 'Filter Tanggal';
  static const String startDate = 'Tanggal Mulai';
  static const String endDate = 'Tanggal Akhir';

  // ============================================
  // Grades
  // ============================================

  static const String gradeHistory = 'Riwayat Nilai';
  static const String gradeSummary = 'Ringkasan Nilai';
  static const String score = 'Nilai';
  static const String average = 'Rata-rata';
  static const String highest = 'Tertinggi';
  static const String lowest = 'Terendah';
  static const String totalGrades = 'Total Nilai';
  static const String teacher = 'Guru';
  static const String subject = 'Mata Pelajaran';
  static const String description = 'Deskripsi';

  // ============================================
  // Homeroom Notes
  // ============================================

  static const String homeroomNotes = 'Catatan Wali Kelas';
  static const String noteFrom = 'Catatan dari';
  static const String noNotes = 'Belum ada catatan';

  // ============================================
  // BK (Bimbingan Konseling)
  // ============================================

  static const String bkInfo = 'Informasi BK';
  static const String achievements = 'Prestasi';
  static const String violations = 'Pelanggaran';
  static const String permits = 'Izin Keluar';
  static const String counseling = 'Konseling';
  static const String counselingNotes = 'Catatan Konseling';
  static const String achievementTitle = 'Judul Prestasi';
  static const String violationCategory = 'Kategori Pelanggaran';
  static const String violationLevel = 'Tingkat';
  static const String violationLevelLight = 'Ringan';
  static const String violationLevelMedium = 'Sedang';
  static const String violationLevelSevere = 'Berat';
  static const String points = 'Poin';
  static const String permitReason = 'Alasan Izin';
  static const String exitTime = 'Waktu Keluar';
  static const String returnTime = 'Waktu Kembali';
  static const String notReturned = 'Belum Kembali';
  static const String downloadDocument = 'Unduh Dokumen';
  static const String parentSummary = 'Ringkasan untuk Orang Tua';

  // ============================================
  // Notifications
  // ============================================

  static const String notificationCenter = 'Pusat Notifikasi';
  static const String markAllRead = 'Tandai Semua Dibaca';
  static const String noNotifications = 'Tidak ada notifikasi';
  static const String unread = 'Belum Dibaca';
  static const String read = 'Sudah Dibaca';
  static const String today = 'Hari Ini';
  static const String yesterday = 'Kemarin';
  static const String thisWeek = 'Minggu Ini';
  static const String earlier = 'Sebelumnya';

  // ============================================
  // Settings
  // ============================================

  static const String accountSettings = 'Pengaturan Akun';
  static const String appSettings = 'Pengaturan Aplikasi';
  static const String notificationSettings = 'Pengaturan Notifikasi';
  static const String about = 'Tentang';
  static const String version = 'Versi';
  static const String privacyPolicy = 'Kebijakan Privasi';
  static const String termsOfService = 'Syarat dan Ketentuan';
  static const String helpSupport = 'Bantuan & Dukungan';

  // ============================================
  // Profile
  // ============================================

  static const String studentProfile = 'Profil Siswa';
  static const String parentProfile = 'Profil Orang Tua';
  static const String nis = 'NIS';
  static const String name = 'Nama';
  static const String className = 'Kelas';
  static const String schoolName = 'Sekolah';
  static const String email = 'Email';
  static const String phone = 'Telepon';

  // ============================================
  // Empty States
  // ============================================

  static const String emptyAttendance = 'Belum ada data kehadiran';
  static const String emptyGrades = 'Belum ada data nilai';
  static const String emptyNotes = 'Belum ada catatan';
  static const String emptyAchievements = 'Belum ada prestasi';
  static const String emptyViolations = 'Tidak ada pelanggaran';
  static const String emptyPermits = 'Tidak ada izin keluar';
  static const String emptyCounseling = 'Tidak ada catatan konseling';
  static const String emptyNotifications = 'Tidak ada notifikasi';
  static const String emptyData = 'Tidak ada data';

  // ============================================
  // Error Messages
  // ============================================

  static const String errorGeneric = 'Terjadi kesalahan. Silakan coba lagi.';
  static const String errorNetwork = 'Tidak dapat terhubung ke server. Periksa koneksi internet Anda.';
  static const String errorTimeout = 'Koneksi timeout. Silakan coba lagi.';
  static const String errorInvalidCredentials = 'NISN atau password salah';
  static const String errorAccountInactive = 'Akun Anda tidak aktif. Hubungi admin sekolah.';
  static const String errorSchoolInactive = 'Sekolah tidak aktif. Hubungi admin.';
  static const String errorSessionExpired = 'Sesi Anda telah berakhir. Silakan login kembali.';
  static const String errorUnauthorized = 'Anda tidak memiliki akses ke halaman ini.';
  static const String errorNotFound = 'Data tidak ditemukan.';
  static const String errorServerError = 'Terjadi kesalahan pada server. Silakan coba lagi nanti.';
  static const String errorPasswordMismatch = 'Password lama tidak sesuai';
  static const String errorSamePassword = 'Password baru tidak boleh sama dengan password lama';
  static const String errorChildNotLinked = 'Anak tidak terhubung dengan akun Anda';

  // ============================================
  // Validation Messages
  // ============================================

  static const String validationRequired = 'Field ini wajib diisi';
  static const String validationNisnRequired = 'NISN wajib diisi';
  static const String validationPasswordRequired = 'Password wajib diisi';
  static const String validationPasswordMinLength = 'Password minimal 6 karakter';
  static const String validationPasswordMismatch = 'Konfirmasi password tidak sesuai';
  static const String validationInvalidFormat = 'Format tidak valid';

  // ============================================
  // Offline Mode
  // ============================================

  static const String offline = 'Offline';
  static const String offlineMode = 'Mode Offline';
  static const String offlineMessage = 'Anda sedang offline. Menampilkan data tersimpan.';
  static const String lastUpdated = 'Terakhir diperbarui';

  // ============================================
  // Date/Time Formats
  // ============================================

  static const String dateFormat = 'dd MMMM yyyy';
  static const String timeFormat = 'HH:mm';
  static const String dateTimeFormat = 'dd MMMM yyyy, HH:mm';
  static const String shortDateFormat = 'dd/MM/yyyy';

  // ============================================
  // Notification Types
  // ============================================

  static const String notificationAttendanceIn = 'Kehadiran Masuk';
  static const String notificationAttendanceOut = 'Kehadiran Pulang';
  static const String notificationViolation = 'Pelanggaran';
  static const String notificationAchievement = 'Prestasi';
  static const String notificationPermit = 'Izin Keluar';
  static const String notificationGrade = 'Nilai Baru';
  static const String notificationHomeroomNote = 'Catatan Wali Kelas';
  static const String notificationCounseling = 'Catatan Konseling';

  // ============================================
  // Helper Methods
  // ============================================

  /// Returns the appropriate greeting based on the hour of day.
  static String getGreeting(int hour) {
    if (hour < 11) return goodMorning;
    if (hour < 15) return goodAfternoon;
    if (hour < 18) return goodEvening;
    return goodNight;
  }

  /// Returns the Indonesian name for an attendance status.
  static String getAttendanceStatusName(String status) {
    switch (status.toLowerCase()) {
      case 'on_time':
        return onTime;
      case 'late':
        return late;
      case 'very_late':
        return veryLate;
      case 'absent':
        return absent;
      case 'sick':
        return sick;
      case 'excused':
        return excused;
      default:
        return status;
    }
  }

  /// Returns the Indonesian name for a violation level.
  static String getViolationLevelName(String level) {
    switch (level.toLowerCase()) {
      case 'ringan':
      case 'light':
        return violationLevelLight;
      case 'sedang':
      case 'medium':
        return violationLevelMedium;
      case 'berat':
      case 'severe':
        return violationLevelSevere;
      default:
        return level;
    }
  }

  /// Returns the Indonesian name for a notification type.
  static String getNotificationTypeName(String type) {
    switch (type.toLowerCase()) {
      case 'attendance_in':
        return notificationAttendanceIn;
      case 'attendance_out':
        return notificationAttendanceOut;
      case 'violation':
        return notificationViolation;
      case 'achievement':
        return notificationAchievement;
      case 'permit':
        return notificationPermit;
      case 'grade':
        return notificationGrade;
      case 'homeroom_note':
        return notificationHomeroomNote;
      case 'counseling':
        return notificationCounseling;
      default:
        return notifications;
    }
  }
}
