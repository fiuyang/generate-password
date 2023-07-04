const fs = require('fs');
const bcrypt = require('bcrypt');
const XLSX = require('xlsx');

// Fungsi untuk mengenkripsi password
function hashPassword(password) {
  const salt = bcrypt.genSaltSync(12);
  const hashedPassword = bcrypt.hashSync(password, salt);
  return hashedPassword;
}

// Fungsi untuk memproses setiap baris pada file Excel
function processRow(row) {
  const updatedRow = { ...row };
  updatedRow.Password = hashPassword(row.Password);
  console.log("Password:", updatedRow.Password);
  return updatedRow;
}

// Fungsi utama
function main() {
  const workbook = XLSX.readFile('Dealer.xlsx');

  const worksheet = workbook.Sheets[workbook.SheetNames[0]];
  const jsonData = XLSX.utils.sheet_to_json(worksheet);

  const processedData = jsonData.map(processRow);

  const newWorkbook = XLSX.utils.book_new();
  const newWorksheet = XLSX.utils.json_to_sheet(processedData);
  XLSX.utils.book_append_sheet(newWorkbook, newWorksheet, 'Sheet1');
  XLSX.writeFile(newWorkbook, 'Dealer.xlsx');
}

main();
