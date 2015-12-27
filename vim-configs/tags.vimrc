
set tags=tags;
set autochdir

" Alias tag selector as C-L
nmap <C-S-n> :ts <c-r>=expand("<cword>")<cr><cr>
