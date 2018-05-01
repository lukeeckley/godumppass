# godumppass
Password dumper for Chrome. Windows only.

## Compilation Tips
You need gcc installed to install go-sqlite3. Cygwin isn't supported. You need to install MinGW. I used the 64bit installer from msys2.org. After installing Msys2, open up the 64bit console and install gcc;

`$ pacman -S mingw-w64-x86_64-gcc`

I temporarily added "C:\msys64\mingw64\bin\" to my %PATH% and was able to install go-sqlite3 afterwards.
