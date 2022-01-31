@cd/d"%~dp0"&(cacls "%SystemDrive%\System Volume Information" >nul 2>nul)||(start "" mshta vbscript:CreateObject^("Shell.Application"^).ShellExecute^("%~nx0"," %*","","runas",1^)^(window.close^)&exit /b)
@echo off
.\main.exe 
pause
