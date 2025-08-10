#!/usr/bin/env python3
"""
Merge split OpenAPI specification files into a single file for oapi-codegen.
"""

import yaml
import os
import sys
from pathlib import Path

def merge_openapi_files():
    """Merge all OpenAPI files into a single specification."""
    api_dir = Path(__file__).parent
    
    # Start with basic structure
    spec = {
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
    schemas_dir = api_dir / "schemas"
    schema_files = ['common.yaml', 'auth.yaml', 'user.yaml', 'vault.yaml', 'audit.yaml', 'api-key.yaml']
    
    for filename in schema_files:
        schema_file = schemas_dir / filename
        if schema_file.exists():
            with open(schema_file, 'r', encoding='utf-8') as f:
                schemas = yaml.safe_load(f)
                if schemas:
                    spec['components']['schemas'].update(schemas)
    
    # Load and merge paths  
    paths_dir = api_dir / "paths"
    
    # Health
    health_file = paths_dir / "health.yaml"
    if health_file.exists():
        with open(health_file, 'r', encoding='utf-8') as f:
            paths = yaml.safe_load(f)
            if 'health' in paths:
                spec['paths']['/api/health'] = paths['health']
    
    # Auth
    auth_file = paths_dir / "auth.yaml"
    if auth_file.exists():
        with open(auth_file, 'r', encoding='utf-8') as f:
            paths = yaml.safe_load(f)
            if 'login' in paths:
                spec['paths']['/api/auth/login'] = paths['login']
            if 'signup' in paths:
                spec['paths']['/api/auth/signup'] = paths['signup']
            if 'logout' in paths:
                spec['paths']['/api/auth/logout'] = paths['logout']
    
    # User
    user_file = paths_dir / "user.yaml"
    if user_file.exists():
        with open(user_file, 'r', encoding='utf-8') as f:
            paths = yaml.safe_load(f)
            if 'getCurrentUser' in paths:
                spec['paths']['/api/user'] = paths['getCurrentUser']
    
    # Vault
    vault_file = paths_dir / "vault.yaml"
    if vault_file.exists():
        with open(vault_file, 'r', encoding='utf-8') as f:
            paths = yaml.safe_load(f)
            if 'vaults' in paths:
                spec['paths']['/api/vaults'] = paths['vaults']
            if 'vault' in paths:
                spec['paths']['/api/vaults/{uniqueId}'] = paths['vault']
    
    # Audit
    audit_file = paths_dir / "audit.yaml"
    if audit_file.exists():
        with open(audit_file, 'r', encoding='utf-8') as f:
            paths = yaml.safe_load(f)
            if 'auditLogs' in paths:
                spec['paths']['/api/audit-logs'] = paths['auditLogs']
    
    # API Key
    apikey_file = paths_dir / "api-key.yaml"
    if apikey_file.exists():
        with open(apikey_file, 'r', encoding='utf-8') as f:
            paths = yaml.safe_load(f)
            if 'apiKeys' in paths:
                spec['paths']['/api/api-keys'] = paths['apiKeys']
            if 'apiKey' in paths:
                spec['paths']['/api/api-keys/{id}'] = paths['apiKey']
    
    # Fix internal references within the merged spec
    def fix_refs(obj):
        if isinstance(obj, dict):
            for key, value in obj.items():
                if key == '$ref' and isinstance(value, str):
                    # Convert all refs to internal refs
                    if value.startswith('../schemas/') or value.startswith('../'):
                        # Extract schema name from ref
                        if '#/' in value:
                            schema_name = value.split('#/')[-1]
                            obj[key] = f'#/components/schemas/{schema_name}'
                    elif value.startswith('#/'):
                        # Already internal refs but may need fixing
                        if not value.startswith('#/components/schemas/'):
                            schema_name = value.replace('#/', '')
                            obj[key] = f'#/components/schemas/{schema_name}'
                else:
                    fix_refs(value)
        elif isinstance(obj, list):
            for item in obj:
                fix_refs(item)
    
    fix_refs(spec)
    
    # Write the merged file
    output_file = api_dir / "api.yaml"
    with open(output_file, 'w', encoding='utf-8') as f:
        yaml.dump(spec, f, default_flow_style=False, allow_unicode=True, sort_keys=False)
    
    print(f"Merged OpenAPI specification written to {output_file}")

if __name__ == "__main__":
    merge_openapi_files()