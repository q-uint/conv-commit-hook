# Git Hook: `commit-msg` for Conventional Commit Messages

## Building

```sh
go build -o pre-commit github.com/q-uint/conv-commit-hook/branch
# cp pre-commit ./.git/hooks/
go build -o commit-msg github.com/q-uint/conv-commit-hook/commit
# cp commit-msg ./.git/hooks/
```

## Installing Hooks Globally

We can make use of Git's [Template Directory](https://git-scm.com/docs/git-init#_template_directory), git will automatically copy files to `$GIT_DIR` after creation.

```sh
mkdir -p ~/.git_template/hooks
git config --global init.templatedir '~/.git_template'
# cp commit-msg ~/.git_template/hooks/
# cp pre-commit ~/.git_template/hooks/
```

If a repository is already initialized, you can just use `git init`.

## Specification Rules

### 01. Commits MUST be prefixed with a type, which consists of a noun, feat, fix, etc., followed by the OPTIONAL scope, OPTIONAL !, and REQUIRED terminal colon and space.

> Supported

Possible Extensions (TODO): fixed list of allowed types.

### 02. The type feat MUST be used when a commit adds a new feature to your application or library.

> NOT Supported: Semantics

### 03. The type fix MUST be used when a commit represents a bug fix for your application.

> NOT Supported: Semantics

### 04. A scope MAY be provided after a type. A scope MUST consist of a noun describing a section of the codebase surrounded by parenthesis, e.g., fix(parser):

> Supported

Possible Extensions (TODO): fixed list of allowed scopes.

### 05. A description MUST immediately follow the colon and space after the type/scope prefix. The description is a short summary of the code changes, e.g., fix: array parsing issue when multiple spaces were contained in string.

> Supported

### 06. A longer commit body MAY be provided after the short description, providing additional contextual information about the code changes. The body MUST begin one blank line after the description.

> Partly Supported (not checked explicitly, we just stop checking stuff after the description)

Possible Extensions (TODO): make this non-optional, for better commit messages?

### 07. A commit body is free-form and MAY consist of any number of newline separated paragraphs.

> TODO?

### 08. One or more footers MAY be provided one blank line after the body. Each footer MUST consist of a word token, followed by either a :<space> or <space># separator, followed by a string value (this is inspired by the git trailer convention).

> TODO?

### 09. A footer’s token MUST use - in place of whitespace characters, e.g., Acked-by (this helps differentiate the footer section from a multi-paragraph body). An exception is made for BREAKING CHANGE, which MAY also be used as a token.

> TODO?

### 10. A footer’s value MAY contain spaces and newlines, and parsing MUST terminate when the next valid footer token/separator pair is observed.

> TODO?

### 11. Breaking changes MUST be indicated in the type/scope prefix of a commit, or as an entry in the footer.

> TODO?

### 12. If included as a footer, a breaking change MUST consist of the uppercase text BREAKING CHANGE, followed by a colon, space, and description, e.g., BREAKING CHANGE: environment variables now take precedence over config files.

> TODO?

### 13. If included in the type/scope prefix, breaking changes MUST be indicated by a ! immediately before the :. If ! is used, BREAKING CHANGE: MAY be omitted from the footer section, and the commit description SHALL be used to describe the breaking change.

> TODO?

### 14. Types other than feat and fix MAY be used in your commit messages, e.g., docs: update ref docs.

> TODO?

### 15. The units of information that make up Conventional Commits MUST NOT be treated as case sensitive by implementors, with the exception of BREAKING CHANGE which MUST be uppercase.

> TODO?

### 16. BREAKING-CHANGE MUST be synonymous with BREAKING CHANGE, when used as a token in a footer.

> TODO?

## References

- Conventional Commits, available under the CC BY 3.0 license, [conventionalcommits.org](https://www.conventionalcommits.org/en/v1.0.0/#specification).
- Conventional Branch, available under the CC BY 4.0 license, [conventional-branch.github.io](https://conventional-branch.github.io/#specification).
- [Git Hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks).

## License

conv-commit-hook is licensed under the GNU Lesser General Public License. See COPYING.LESSER for details.
