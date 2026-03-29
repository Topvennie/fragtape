Start-Sleep -Seconds 10

& "$env:SystemRoot\System32\DisplaySwitch.exe" /external

Start-Sleep -Seconds 2

Start-Process "powershell.exe" -ArgumentList "-NoProfile -ExecutionPolicy Bypass -File C:\Fragtape\scripts\watch.ps1"
