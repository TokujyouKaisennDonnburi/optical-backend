@echo off
rem Load environment variables from .env file
for /f "tokens=1,2 delims==" %%a in ('type .env ^| findstr /v "^#" ^| findstr /v "^$"') do set %%a=%%b
sql-migrate up