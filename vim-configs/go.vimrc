" @require: github.com/fatih/vim-go
"

" set for go
let g:go_highlight_functions = 1
let g:go_highlight_methods = 1
let g:go_highlight_structs = 1
let g:go_highlight_operators = 1
let g:go_highlight_build_constraints = 1

" go auto-import
let g:go_fmt_command = "goimports"

" shows errors for the fmt command, to disable it:
let g:go_fmt_fail_silently = 1

" Disable auto fmt on save
let g:go_fmt_autosave = 1

" Disable opening browser after posting your snippet to play.golang.org
let g:go_play_open_browser = 0

" Go related mappings
au FileType go nmap <Leader>r <Plug>(go-run)
au FileType go nmap <Leader>b <Plug>(go-build)
au FileType go nmap <Leader>t <Plug>(go-test)
au FileType go nmap <leader>c <Plug>(go-coverage)
au FileType go nmap <Leader>i <Plug>(go-info)

" Rename the identifier under the cursor to a new name
au FileType go nmap <Leader>e <Plug>(go-rename)

" Open the relevant Godoc for the word under the cursor
au FileType go nmap <Leader>gd <Plug>(go-doc)
au FileType go nmap <Leader>gv <Plug>(go-doc-vertical)

au FileType go nmap gd <Plug>(go-def-tab)

" open the definition/declaration, in a new vertical, horizontal, or tab, for the word under your cursor
au FileType go nmap <Leader>ds <Plug>(go-def-split)
au FileType go nmap <Leader>dv <Plug>(go-def-vertical)
au FileType go nmap <Leader>dt <Plug>(go-def-tab)

" open the Godoc in browser
au FileType go nmap <Leader>gb <Plug>(go-doc-browser)

" Show a list of interfaces which is implemented by the type under your cursor
au FileType go nmap <Leader>s <Plug>(go-implements)

" Show type info for the word under your cursor
au FileType go nmap <Leader>i <Plug>(go-info)
