@echo off
setlocal

tzutil /s "UTC"

mkdir C:\Fragtape\scripts 2>nul
mkdir C:\Fragtape\logs 2>nul

copy /Y C:\OEM\watch.ps1 C:\Fragtape\scripts\watch.ps1 >nul

schtasks /Delete /F /TN "FragtapeRecorderDev" >nul 2>&1

schtasks /Create /F ^
  /TN "FragtapeRecorderDev" ^
  /SC ONSTART ^
  /RU "SYSTEM" ^
  /RL HIGHEST ^
  /TR "powershell.exe -NoProfile -ExecutionPolicy Bypass -File C:\Fragtape\scripts\watch.ps1" >nul

schtasks /Run /TN "FragtapeRecorderDev" >nul 2>&1

endlocal

