#!/bin/bash

# Output file name
OUTPUT_FILE="web_page_analyser_project.txt"

# Function to create separator
create_separator() {
    local filename=$1
    printf "\n=====\n"
    printf "%s\n" "$filename"
    printf "=======\n"
}

# Function to check if file is binary
is_binary() {
    local file=$1
    if file "$file" | grep -q "text"; then
        return 1  # Not binary
    else
        return 0  # Binary
    fi
}

# Clear or create the output file
> "$OUTPUT_FILE"

# Find all files, excluding .git, node_modules, and other unnecessary directories
find . -type f \
    ! -path "./.git/*" \
    ! -path "./web/node_modules*" \
    ! -path "./archive/*" \
    ! -path "./db_migration/*" \
    ! -path "./sql/*" \
    ! -path "./__pycache__/*" \
    ! -path "./coverage/*" \
    ! -path "./*.md" \
    ! -path "./*.sh" \
    ! -path "./*.json" \
    ! -path "./*.sum" \
    ! -path "./*.json.back" \
    ! -path "./.gitignore" \
    ! -path "./$OUTPUT_FILE" \
    ! -name "package-lock.json" \
    ! -name "ecs-definition-web-live.tpl" \
    ! -name "jest.config.js" \
    ! -name "DOCUMENTATION.md" \
    ! -name "SWAGGER_DEVELOPER.md" \
    ! -name "SWAGGER_INTEGRATION.md" \
    ! -name "SWAGGER_README.md" \
    ! -name "test.config.js" \
    ! -name "README.md" \
    ! -name "concatenated_project.txt" \
    ! -name "condens.sh" \
    -print0 | while IFS= read -r -d '' file; do
    
    # Skip binary files
    if is_binary "$file"; then
        echo "Skipping binary file: $file"
        continue
    fi

    # Get relative path
    relative_path=${file#./}
    
    # Create separator with filename
    create_separator "$relative_path" >> "$OUTPUT_FILE"
    
    # Append file contents
    cat "$file" >> "$OUTPUT_FILE"
    
    echo "Processed: $relative_path"
done

echo "Concatenation complete! Output written to: $OUTPUT_FILE"
ls -l "$OUTPUT_FILE"