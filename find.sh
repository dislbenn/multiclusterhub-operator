#!/bin/bash

# Define base directories
CHARTS_DIR="pkg/templates/charts"
CRDS_DIR="pkg/templates/crds"
OUTPUT_FILE="report.md"

# Write Markdown table header
echo "| File Path | Kind / CRD Name |" > "$OUTPUT_FILE"
echo "|-----------|----------------|" >> "$OUTPUT_FILE"

# Process general YAML files in CHARTS_DIR
find "$CHARTS_DIR" -type f \( -name "*.yaml" -o -name "*.yml" \) | while read -r file; do
    # Extract all `kind` values in case of multiple documents
    yq e '.kind' "$file" 2>/dev/null | while read -r kind; do
        if [[ -n "$kind" && "$kind" != "null" ]]; then
            echo "| $file | $kind |" >> "$OUTPUT_FILE"
        fi
    done
done

# Process CRD YAML files in CRDS_DIR
find "$CRDS_DIR" -type f \( -name "*.yaml" -o -name "*.yml" \) | while read -r file; do
    # Extract all `kind` values in case of multiple documents
    yq e '.kind' "$file" 2>/dev/null | while read -r kind; do
        if [[ "$kind" == "CustomResourceDefinition" ]]; then
            crd_name=$(yq e '.metadata.name' "$file" 2>/dev/null)
            if [[ -n "$crd_name" && "$crd_name" != "null" ]]; then
                echo "| $file | $crd_name |" >> "$OUTPUT_FILE"
            fi
        fi
    done
done

echo "Markdown report saved to $OUTPUT_FILE"
