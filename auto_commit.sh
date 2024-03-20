#!/bin/bash

# Função para gerar a mensagem de commit com base no tipo de alteração
generate_commit_message() {
    local filename="$1"
    local message=""

    # Extrair o tipo de alteração do nome do arquivo
    local change_type=$(echo "$filename" | cut -d '/' -f 1)

    case "$change_type" in
        "backend")
            message="chore(backend): Alterações em $filename"
            ;;
        "frontend")
            message="chore(frontend): Alterações em $filename"
            ;;
        *)
            message="chore: Alterações em $filename"
            ;;
    esac

    echo "$message"
}

# Exibir o número de arquivos modificados
echo "Verificando número de arquivos modificados..."
num_modified_files=$(git status -s | wc -l)
echo "Número de arquivos modificados: $num_modified_files"

# Verificar se não há arquivos modificados
if [ $num_modified_files -eq 0 ]; then
    echo "Nenhum arquivo modificado. Interrompendo a execução do script."
    exit 0
fi

# Adicionar todos os arquivos modificados ao commit
echo "Adicionando todos os arquivos modificados ao commit..."
git add -A

# Realizar commit
echo "Realizando commit..."
git commit -m "chore: Atualizações automáticas"

# Exibir todos os commits gerados e adicionados
echo "Todos os commits gerados e adicionados:"
git log --oneline

# Verificar se os commits foram feitos
if [ $? -eq 0 ]; then
    # Enviar alterações para o repositório remoto
    echo "Enviando alterações para o repositório remoto..."
    git push origin main
    echo "Alterações enviadas com sucesso!"
else
    echo "Ocorreu um erro ao realizar os commits. Verifique o status do git."
fi
