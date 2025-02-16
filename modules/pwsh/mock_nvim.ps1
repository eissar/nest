return @{
    Details = @{
        COMMAND = 'wt.exe'
        VIM = 'D:\Neovim\v0.4.2'
        VIMRUNTIME = 'D:\Neovim\v0.4.2\share\nvim\runtime'
        DESCRIPTION = 'open neovim v0.3.0 in new terminal tab'
        ISBACKGROUND = $true
    }
} | ConvertTo-Json -Depth 7 -AsArray
