@echo off
setlocal

mkdir C:\Fragtape 2>nul
mkdir C:\Fragtape\config 2>nul
mkdir C:\Fragtape\logs 2>nul

copy /Y C:\OEM\recorder.exe C:\Fragtape\recorder.exe
copy /Y C:\OEM\config\production.yml C:\Fragtape\config\production.yml

tzutil /s "UTC"

schtasks /Create /F /RL HIGHEST /SC ONSTART /TN "FragtapeRecorder" ^
  /TR "cmd /c C:\Fragtape\recorder.exe >> C:\Fragtape\logs\stdout.log 2>> C:\Fragtape\logs\stderr.log"

schtasks /Run /TN "FragtapeRecorder"

endlocal
