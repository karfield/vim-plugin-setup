
set textwidth=80
set scrolloff=7
set sidescroll=2

set nonu

set laststatus=2
set showcmd		    " Show (partial) command in status line.

set incsearch

set autoindent
set si "smartident
set cindent
"set cinoptions=g0,N-s
set cinoptions=h0,l0,g0,t0,i4,+4,(0,w1,W4

set smarttab
set tabstop=4
set softtabstop=4
set shiftwidth=4
set expandtab

hi Normal ctermbg=none	"Transparent background
"set background=dark

colorscheme koehler

set nowrap	"disable wrap
set linebreak

set showmatch		" Show matching brackets.
set matchtime=1		" match time, 1s

set backspace=indent,eol,start	"disable Backspace/Delete when in visual

"remap some match pairs for editing
inoremap [] []<Left>
inoremap () ()<Left>
inoremap {} {}<Left>
inoremap "" ""<Left>
inoremap <> <><Left>

