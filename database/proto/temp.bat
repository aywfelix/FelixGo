@echo off
for %%i in (*.proto) do (
    protoc  --go_out=. %%i
)

pause