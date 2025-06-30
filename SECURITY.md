# Security Policy

## 🔒 Credential Protection

This repository has multiple layers of protection against accidental credential commits:

### 1. Pre-commit Hooks
- Automatically scans files before commit for potential secrets
- Checks for common patterns: passwords, API keys, tokens
- Prevents commits of dangerous files like `.zmcp.json`, `.env`
- Located in `.githooks/pre-commit`

To enable: `git config core.hooksPath .githooks`

### 2. .gitignore
Explicitly ignores:
- `*.zmcp.json` - MCP configuration files
- `*.key`, `*.pem` - Private keys
- `.env*` - Environment files
- `secrets/`, `credentials/` - Secret directories

### 3. GitHub Actions Security Scanning
- **Gitleaks**: Scans for secrets in code and history
- **Gosec**: Go-specific security analysis
- **Dependency scanning**: Checks for vulnerable dependencies

### 4. Gitleaks Configuration
- Custom patterns in `.gitleaks.toml`
- Detects MCP configs, API keys, tokens, passwords
- Allows test/example passwords

## 🚨 If You've Committed Secrets

1. **Immediately rotate** the exposed credentials
2. **Remove from history**:
   ```bash
   # Remove file from all history
   git filter-branch --force --index-filter \
     'git rm --cached --ignore-unmatch PATH_TO_FILE' \
     --prune-empty --tag-name-filter cat -- --all
   
   # Force push (coordinate with team)
   git push origin --force --all
   git push origin --force --tags
   ```

3. **Contact security team** if credentials were exposed publicly

## 🛡️ Best Practices

1. **Never commit**:
   - Real passwords, API keys, or tokens
   - `.zmcp.json` or similar config files
   - Private keys or certificates
   - `.env` files with real values

2. **Use instead**:
   - Environment variables
   - External secret management
   - Example/template files (e.g., `.env.example`)
   - Placeholder values in docs

3. **Before committing**:
   - Review `git diff` carefully
   - Check `git status` for unexpected files
   - Run `gitleaks detect --staged`

## 🔍 Manual Security Scan

```bash
# Install gitleaks
brew install gitleaks

# Scan current directory
gitleaks detect

# Scan staged changes
gitleaks detect --staged

# Scan with custom config
gitleaks detect -c .gitleaks.toml
```

## 📝 Reporting Security Issues

If you discover a security vulnerability:

1. **Do NOT** create a public issue
2. Email security concerns to [maintainer email]
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

## 🔄 Regular Audits

We perform regular security audits:
- Weekly automated scans via GitHub Actions
- Dependency updates monthly
- Manual review quarterly

## 📚 Resources

- [OWASP Go Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Go_Security_Cheat_Sheet.html)
- [GitHub Secret Scanning](https://docs.github.com/en/code-security/secret-scanning)
- [Gitleaks Documentation](https://github.com/gitleaks/gitleaks)