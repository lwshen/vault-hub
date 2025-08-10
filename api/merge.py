#!/usr/bin/env python3
"""
Script to merge multiple OpenAPI YAML files into a single file.
This allows us to split the API definition into multiple files while
still generating a single file for oapi-codegen.
"""

import yaml
import os
import sys
from pathlib import Path

def load_yaml(file_path):
    """Load a YAML file and return its content."""
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            return yaml.safe_load(f)
    except Exception as e:
        print(f"Error loading {file_path}: {e}")
        return None

def merge_schemas(schema_files):
    """Merge schema files into a single components.schemas section."""
    merged_schemas = {}
    
    for schema_file in schema_files:
        if os.path.exists(schema_file):
            content = load_yaml(schema_file)
            if content:
                # Each schema file contains direct schema definitions
                for schema_name, schema_def in content.items():
                    merged_schemas[schema_name] = schema_def
    
    return merged_schemas

def resolve_refs_in_paths(paths, schemas):
    """Resolve $ref references in paths to use inline schemas."""
    import copy
    resolved_paths = copy.deepcopy(paths)
    
    def resolve_ref(obj):
        if isinstance(obj, dict):
            for key, value in obj.items():
                if key == '$ref' and isinstance(value, str):
                    # Extract schema name from ref
                    if '#/' in value:
                        schema_name = value.split('#/')[-1]
                        if schema_name in schemas:
                            return schemas[schema_name]
                elif isinstance(value, (dict, list)):
                    obj[key] = resolve_ref(value)
        elif isinstance(obj, list):
            for i, item in enumerate(obj):
                if isinstance(item, (dict, list)):
                    obj[i] = resolve_ref(item)
        return obj
    
    for path, path_def in resolved_paths.items():
        resolved_paths[path] = resolve_ref(path_def)
    
    return resolved_paths

def resolve_refs_in_schemas(schemas):
    """Resolve $ref references within schemas to use inline schemas."""
    import copy
    resolved_schemas = copy.deepcopy(schemas)
    
    def resolve_ref(obj):
        if isinstance(obj, dict):
            for key, value in obj.items():
                if key == '$ref' and isinstance(value, str):
                    # Extract schema name from ref
                    if '#/' in value:
                        schema_name = value.split('#/')[-1]
                        if schema_name in resolved_schemas:
                            return resolved_schemas[schema_name]
                elif isinstance(value, (dict, list)):
                    obj[key] = resolve_ref(value)
        elif isinstance(obj, list):
            for i, item in enumerate(obj):
                if isinstance(item, (dict, list)):
                    obj[i] = resolve_ref(item)
        return obj
    
    for schema_name, schema_def in resolved_schemas.items():
        resolved_schemas[schema_name] = resolve_ref(schema_def)
    
    return resolved_schemas

def main():
    # Define the base OpenAPI structure
    base_openapi = {
        'openapi': '3.0.0',
        'info': {
            'version': '1.0.0',
            'title': 'Vault Hub Server'
        },
        'paths': {},
        'components': {
            'schemas': {}
        }
    }
    
    # Load and merge schemas
    schema_dir = Path('schemas')
    schema_files = [
        schema_dir / 'common.yaml',
        schema_dir / 'auth.yaml',
        schema_dir / 'user.yaml',
        schema_dir / 'vault.yaml',
        schema_dir / 'api_key.yaml',
        schema_dir / 'audit_log.yaml'
    ]
    
    base_openapi['components']['schemas'] = merge_schemas(schema_files)
    
    # Load and merge paths
    paths_dir = Path('paths')
    path_files = [
        paths_dir / 'health.yaml',
        paths_dir / 'auth.yaml',
        paths_dir / 'user.yaml',
        paths_dir / 'vaults.yaml',
        paths_dir / 'api_keys.yaml',
        paths_dir / 'audit_logs.yaml'
    ]
    
    for path_file in path_files:
        if os.path.exists(path_file):
            content = load_yaml(path_file)
            if content:
                # Each path file contains path definitions
                for path, path_def in content.items():
                    base_openapi['paths'][path] = path_def
    
    # First resolve references within schemas
    base_openapi['components']['schemas'] = resolve_refs_in_schemas(
        base_openapi['components']['schemas']
    )
    
    # Then resolve all $ref references in paths to use inline schemas
    base_openapi['paths'] = resolve_refs_in_paths(
        base_openapi['paths'], 
        base_openapi['components']['schemas']
    )
    
    # Write the merged file
    output_file = 'merged_api.yaml'
    with open(output_file, 'w', encoding='utf-8') as f:
        yaml.dump(base_openapi, f, default_flow_style=False, sort_keys=False, allow_unicode=True)
    
    print(f"Successfully merged OpenAPI files into {output_file}")
    return 0

if __name__ == '__main__':
    sys.exit(main())