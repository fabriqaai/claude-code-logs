---
name: code-review
type: knowledge
version: 1.0.0
agent: CodeActAgent
triggers:
  - any
  - commit
---

# Code Review Microagent

This microagent provides guidance for performing comprehensive code reviews on commits and changes in the repository.

## Purpose

Automatically triggered when reviewing any code changes or commits to ensure consistent, thorough code review practices.

## Code Review Checklist

When performing a code review, evaluate the following aspects:

### 1. Code Quality
- Is the code clean, readable, and well-organized?
- Are variable and function names descriptive and meaningful?
- Is there unnecessary code duplication that could be refactored?
- Are functions and methods appropriately sized and focused?

### 2. Logic and Correctness
- Does the code correctly implement the intended functionality?
- Are edge cases handled appropriately?
- Is error handling implemented correctly?
- Are there any potential bugs or logic errors?

### 3. Performance
- Are there any obvious performance issues?
- Are loops and iterations efficient?
- Is memory usage reasonable?
- Are there unnecessary computations or redundant operations?

### 4. Security
- Are there any security vulnerabilities?
- Is user input properly validated and sanitized?
- Are sensitive data handled securely?
- Are there any potential injection vulnerabilities?

### 5. Testing
- Are there adequate tests for the changes?
- Do the tests cover edge cases?
- Are the tests meaningful and not just for coverage?

### 6. Documentation
- Is the code adequately commented where necessary?
- Are complex algorithms or business logic explained?
- Is the README or other documentation updated if needed?

### 7. Style and Conventions
- Does the code follow the project's coding standards?
- Is formatting consistent with the rest of the codebase?
- Are imports organized properly?

## Review Process

1. **Understand the context**: Read the commit message and any related issues
2. **Review the diff**: Examine each changed file systematically
3. **Test locally if needed**: For significant changes, consider running the code
4. **Provide constructive feedback**: Be specific, actionable, and respectful
5. **Approve or request changes**: Make a clear decision with justification

## Feedback Guidelines

- Be specific about what needs to change and why
- Suggest improvements rather than just pointing out problems
- Acknowledge good practices and well-written code
- Prioritize feedback by importance (blocking vs. nice-to-have)
