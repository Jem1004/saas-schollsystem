import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'

export interface ExportColumn {
  header: string
  dataKey: string
  width?: number
}

export interface ExportOptions {
  title: string
  subtitle?: string
  filename: string
  columns: ExportColumn[]
  data: Record<string, unknown>[]
  schoolName?: string
  dateRange?: { start: string; end: string }
}

const formatDate = (date: Date): string => {
  return date.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

export const exportToPDF = (options: ExportOptions): void => {
  const { title, subtitle, filename, columns, data, schoolName, dateRange } = options

  const doc = new jsPDF()
  const pageWidth = doc.internal.pageSize.getWidth()

  // Header
  let yPos = 20

  // School name
  if (schoolName) {
    doc.setFontSize(14)
    doc.setFont('helvetica', 'bold')
    doc.text(schoolName, pageWidth / 2, yPos, { align: 'center' })
    yPos += 8
  }

  // Title
  doc.setFontSize(16)
  doc.setFont('helvetica', 'bold')
  doc.text(title, pageWidth / 2, yPos, { align: 'center' })
  yPos += 7

  // Subtitle
  if (subtitle) {
    doc.setFontSize(10)
    doc.setFont('helvetica', 'normal')
    doc.text(subtitle, pageWidth / 2, yPos, { align: 'center' })
    yPos += 6
  }

  // Date range
  if (dateRange) {
    doc.setFontSize(10)
    doc.setFont('helvetica', 'normal')
    doc.text(`Periode: ${dateRange.start} - ${dateRange.end}`, pageWidth / 2, yPos, { align: 'center' })
    yPos += 6
  }

  // Generated date
  doc.setFontSize(9)
  doc.setFont('helvetica', 'italic')
  doc.text(`Dicetak: ${formatDate(new Date())}`, pageWidth / 2, yPos, { align: 'center' })
  yPos += 10

  // Table
  const tableColumns = columns.map(col => col.header)
  const tableData = data.map(row => columns.map(col => String(row[col.dataKey] ?? '')))

  autoTable(doc, {
    head: [tableColumns],
    body: tableData,
    startY: yPos,
    styles: {
      fontSize: 9,
      cellPadding: 3,
    },
    headStyles: {
      fillColor: [249, 115, 22], // Orange color
      textColor: 255,
      fontStyle: 'bold',
    },
    alternateRowStyles: {
      fillColor: [254, 249, 243],
    },
    margin: { left: 14, right: 14 },
  })

  // Footer with page numbers
  const pageCount = (doc as unknown as { internal: { getNumberOfPages: () => number } }).internal.getNumberOfPages()
  for (let i = 1; i <= pageCount; i++) {
    doc.setPage(i)
    doc.setFontSize(8)
    doc.setFont('helvetica', 'normal')
    doc.text(
      `Halaman ${i} dari ${pageCount}`,
      pageWidth / 2,
      doc.internal.pageSize.getHeight() - 10,
      { align: 'center' }
    )
  }

  doc.save(`${filename}.pdf`)
}

// Helper to format violation data for export
export const formatViolationForExport = (violation: Record<string, unknown>): Record<string, unknown> => {
  const levelLabels: Record<string, string> = {
    ringan: 'Ringan',
    sedang: 'Sedang',
    berat: 'Berat',
  }
  return {
    ...violation,
    createdAt: new Date(violation.createdAt as string).toLocaleDateString('id-ID'),
    level: levelLabels[violation.level as string] || violation.level,
  }
}

// Helper to format achievement data for export
export const formatAchievementForExport = (achievement: Record<string, unknown>): Record<string, unknown> => {
  return {
    ...achievement,
    createdAt: new Date(achievement.createdAt as string).toLocaleDateString('id-ID'),
    point: `+${achievement.point}`,
  }
}

// Helper to format permit data for export
export const formatPermitForExport = (permit: Record<string, unknown>): Record<string, unknown> => {
  const formatTime = (dateStr: string) => {
    return new Date(dateStr).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
  }
  return {
    ...permit,
    createdAt: new Date(permit.createdAt as string).toLocaleDateString('id-ID'),
    exitTime: formatTime(permit.exitTime as string),
    returnTime: permit.returnTime ? formatTime(permit.returnTime as string) : 'Belum Kembali',
    status: permit.returnTime ? 'Sudah Kembali' : 'Belum Kembali',
  }
}

// Helper to format counseling data for export
export const formatCounselingForExport = (note: Record<string, unknown>): Record<string, unknown> => {
  return {
    ...note,
    createdAt: new Date(note.createdAt as string).toLocaleDateString('id-ID'),
    status: note.parentSummary ? 'Dibagikan' : 'Internal',
  }
}


// Interface for permit document data from backend
export interface PermitDocumentData {
  permit_id: number
  student_name: string
  student_nis: string
  student_nisn: string
  class_name: string
  school_name: string
  reason: string
  exit_time: string
  return_time?: string
  responsible_teacher: string
  generated_at: string
}

// Generate permit document PDF in receipt/thermal printer format (80mm width)
export const generatePermitPDF = (data: PermitDocumentData): void => {
  // 80mm width = ~226 points (72 points per inch, 80mm = 3.15 inches)
  const receiptWidth = 226
  const margin = 10
  const contentWidth = receiptWidth - (margin * 2)
  
  // Create PDF with custom receipt size
  const doc = new jsPDF({
    orientation: 'portrait',
    unit: 'pt',
    format: [receiptWidth, 400],
  })

  let yPos = 15
  const lineHeight = 12
  const smallLineHeight = 10
  const centerX = receiptWidth / 2

  // Helper: draw dashed line using text
  const drawDashedLine = (y: number) => {
    doc.setFontSize(6)
    doc.setFont('helvetica', 'normal')
    doc.text('- - - - - - - - - - - - - - - - - - - - - - - - - - - -', centerX, y, { align: 'center' })
  }

  // Helper: print label-value line
  const printLine = (label: string, value: string) => {
    doc.setFontSize(8)
    doc.setFont('helvetica', 'normal')
    doc.text(label, margin, yPos)
    doc.text(':', margin + 55, yPos)
    const maxValueWidth = contentWidth - 60
    const lines = doc.splitTextToSize(value, maxValueWidth)
    doc.text(lines, margin + 60, yPos)
    yPos += smallLineHeight * Math.max(lines.length, 1)
  }

  // === HEADER ===
  doc.setFontSize(10)
  doc.setFont('helvetica', 'bold')
  doc.text(data.school_name, centerX, yPos, { align: 'center', maxWidth: contentWidth })
  yPos += lineHeight + 2

  doc.setFontSize(9)
  doc.setFont('helvetica', 'bold')
  doc.text('SURAT IZIN KELUAR', centerX, yPos, { align: 'center' })
  yPos += lineHeight

  doc.setFontSize(7)
  doc.setFont('helvetica', 'normal')
  doc.text(`No: IK/${data.permit_id}/${new Date().getFullYear()}`, centerX, yPos, { align: 'center' })
  yPos += smallLineHeight + 2

  drawDashedLine(yPos)
  yPos += 8

  // === STUDENT INFO ===
  printLine('Nama', data.student_name)
  printLine('NIS', data.student_nis)
  printLine('Kelas', data.class_name)
  yPos += 3

  drawDashedLine(yPos)
  yPos += 8

  // === PERMIT DETAILS ===
  printLine('Alasan', data.reason)
  printLine('Waktu Keluar', formatReceiptDateTime(data.exit_time))
  printLine('Guru PJ', data.responsible_teacher)
  
  if (data.return_time) {
    printLine('Waktu Kembali', formatReceiptDateTime(data.return_time))
  }
  yPos += 3

  drawDashedLine(yPos)
  yPos += 10

  // === SIGNATURE ===
  doc.setFontSize(7)
  doc.text('Guru BK,', centerX, yPos, { align: 'center' })
  yPos += 25
  doc.text('(_______________)', centerX, yPos, { align: 'center' })
  yPos += 12

  drawDashedLine(yPos)
  yPos += 8

  // === FOOTER ===
  doc.setFontSize(6)
  doc.setFont('helvetica', 'italic')
  doc.text(formatReceiptDateTime(data.generated_at), centerX, yPos, { align: 'center' })
  yPos += smallLineHeight
  doc.text('Simpan struk ini sebagai bukti', centerX, yPos, { align: 'center' })

  // Open in new tab
  const pdfBlob = doc.output('blob')
  const url = URL.createObjectURL(pdfBlob)
  window.open(url, '_blank')
}

// Format datetime for receipt (compact format)
const formatReceiptDateTime = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }) + ' ' + date.toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
  })
}
