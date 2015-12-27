" @require: github.com/bling/vim-airline
"
" @run-script
" #!/bin/bash
" mkdir -p ${VIMDIR}/tmp
" pushd ${VIMDIR}/tmp
" git clone https://github.com/powerline/fonts.git
" pushd fonts
" ./install.sh
" popd
" popd
" @end-script
" 
"
" https://github.com/bling/vim-airline
set guifont=Inconsolata\ for\ Powerline:h18
set encoding=utf-8
"let g:Powerline_symbols = 'fancy'
set t_Co=256
set fillchars+=stl:\ ,stlnc:\
set term=xterm-256color
set termencoding=utf-8

let g:airline_powerline_fonts = 1

if has("gui_running")
    let s:uname = system("uname")
    if s:uname == "Darwin\n"
        set guifont=Inconsolata\ for\ Powerline:h18
    endif
endif
