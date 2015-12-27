
" Copy & Paste
"
" For Mac
vmap <C-c> y:call system("pbcopy", getreg("\""))<CR>
nmap <C-s> :call setreg("\"", system("pbpaste"))<CR>
