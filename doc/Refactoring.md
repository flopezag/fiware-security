# 🛠️ FIWARE Security Scanner – Refactoring Plan

## 🎯 Objectives
- Refactor code for better maintainability and reliability.
- Improve risk calculation using the formula:  
- Add automated workflows via GitHub Actions.
- Build and publish a Docker image.
- Improve CLI UX and JSON-driven logic.
- Update README and usage documentation.
- Review count Risk cases and add number of critical issues
- Integrate Gitleaks and check Docker Bench Security

---

## 📦 Phase 1: Codebase Refactoring

### ✅ Goals
- Modularize code for easier maintenance.
- Update risk calculation logic.
- Ensure clear error messages and logging.
- Review count Risk cases and add number of critical issues.
- Integrate Gitleaks and check Docker Bench Security.
- Unify the file enablers.json
- Automatic email generation with the generated output of the analysis (results data.)
- Create tests

### 🔧 Tasks
- [x] **Refactor `analyzeAndVisualize`**:
- [x] Replace mean risk with: `1 - product(1 - risk_i)`.
  (note if any risk is 0, final probability is 1, no change in the calculation)
- [x] Create **helper functions** for:
    - Loading and parsing `enablers.json`.
    - Matching enabler names to report filenames.
    - Finding the most recent result file.
- [ ] Improve CLI validation (e.g., missing enabler or file).
- [ ] Improve error handling and user-friendly messages.
- [ ] Review the calculation of number of Risk.
- [ ] Add the total number of critical issues.
- [ ] Create some tests files associated to the code.
- [ ] Keep only enablers.json in config file.
- [ ] Auto-email report summary to the email field in enablers.json
 
---

## 🧪 Phase 2: GitHub Actions Integration

### ✅ Goals
- Automate linting, testing, and Docker publishing

### 🔧 Tasks
- [ ] Adding Lint + Test files to `.github/workflows`

---

## 🐳 Phase 3: Docker Support

### ✅ Goals
- Let users run the scanner without installing Go

### 🔧 Tasks
- [ ] Build locally the docker image of the service.
- [ ] Publish the image into fiware organization.
- [ ] Check how to store or save the data generated in the process.

---

## 📘 Phase 4: README Update

### ✅ Goals
- Provide clear instructions and examples

### 🔧 Tasks
- [ ] Add project summary and goals.
- [ ] Describe installation via Go and Docker.
- [ ] Describe usage of the application.
- [ ] Add sample output (with anonymized data).
- [ ] Document enablers.json structure (array format).
- [ ] Explain how to modify the content of enablers.json
- [ ] Add status badges (CI build, Go Report Card, DockerHub availability, License)

---

## 🧠 Phase 5: Optional Enhancements

### ✅ Goals
- Develop future functionalities to the application.

### 🔧 Tasks

- [ ] Generate Markdown or HTML vulnerability reports.
- [ ] Support SARIF or SBOM output formats.
- [ ] Add --severity or --epss-min filters.
