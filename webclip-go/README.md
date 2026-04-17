# Webclip (Go Port)

Webclip is a static website generator specifically designed for creating dictionary-like lexicons or structured tutorials. It transforms raw HTML content into a polished, navigable website with automatic indexing, searching, and menu generation.

This version is a high-performance port to **Go**, using only the standard library.

---

## 📂 Project Structure

A typical Webclip installation looks like this:

```text
webclip-go/
├── webclip             # The compiled executable
├── config.json         # Global configuration
├── templates/          # Directory containing site themes
│   └── default/        # A theme folder
│       ├── config.json # Theme configuration
│       └── index.html  # The main layout wrap
└── projects/           # Directory containing your project sources
    └── my-lexicon/     # A specific project folder
        ├── config.json # Project-specific configuration
        └── src/        # Your raw HTML content files
```

---

## ⚙️ Configuration

### 1. Global Configuration (`config.json`)
Located in the root folder, this tells the tool where to find your projects.
```json
{
  "projectsPath": "projects"
}
```

### 2. Project Configuration (`projects/NAME/config.json`)
Each project has its own settings:
```json
{
  "name": "My Great Lexicon",  // The display name and output folder name
  "template": "default",      // Which folder in /templates to use
  "variables": {              // Global variables available in all templates
    "projectTitle": "Title",
    "projectSubTitle": "Subtitle",
    "topics": ["Intro", "Advanced"] // (Optional) Enables topic-based navigation
  },
  "parse": ["*.html"]         // Glob patterns for files in /src to process
}
```

### 3. Template Configuration (`templates/NAME/config.json`)
Defines how the theme assets are handled:
```json
{
  "parse": ["styles"],        // Folders with variables to replace (e.g., CSS)
  "copy": ["scripts", "fonts"], // Folders to copy directly to output
  "variables": {              // Default theme variables (can be overridden by project)
    "color": "#abcdef"
  }
}
```

---

## 📝 Content Creation & Metadata

Webclip extracts logic from your HTML source files using standard `<meta>` tags. Add these inside the `<head>` of your files in `src/`:

| Meta Attribute | Description |
| :--- | :--- |
| `data-menu` | Assigns the page to a menu (e.g., `main`, `top`). Default is `main`. |
| `data-sort` | The key used for sorting. |
| | - In **Alphabetical** mode: Overrides the title (e.g., "Arithmetisches Mittel" instead of "Mittelwert"). |
| | - In **Topic** mode: Defines which topic group the page belongs to. |
| `data-weight` | A number used for ordering (e.g., `10`, `20`). Lower numbers appear first. |
| `data-alias` | An alternate name for the entry. Useful for "See also" style navigation. |
| `name="keywords"`| Comma-separated list of words added to the `search.json` index. |

**Example:**
```html
<head>
  <title>Mittelwert</title>
  <meta data-menu="main" data-sort="Arithmetisches Mittel" data-alias="Average">
  <meta name="keywords" content="Statistik, Durchschnitt, Mathe">
</head>
```

---

## 🧭 Navigation Styles

Webclip supports two primary navigation modes based on your `variables`:

### 1. Alphabetical (Lexikon Style)
This is the default mode. Pages are grouped by their starting letter.
- **Trigger**: No `topics` defined in `config.json`.
- **Behavior**: Uses the page `<title>` or `data-sort` key to build an A-Z sidebar.

### 2. Topic-Based (Tutorial Style)
Pages are grouped into manual categories.
- **Trigger**: Define a `"topics": ["A", "B", "C"]` array in the project variables.
- **Behavior**: Matches the page's `data-sort` against the defined topics. Pages without a category appear under **"???"**.

---

## 🎨 Templating

Webclip uses a simple `{{ variable }}` syntax for substitution.

1.  **Asset Parsing**: Files in the template's `parse` folders (like CSS) will have variables like `{{ color }}` replaced.
2.  **HTML Wrapping**: Every content file is wrapped inside the template's `index.html`.
    - `{{ body }}` is replaced with everything inside the source file's `<body>` tag.
    - `{{ title }}` is replaced with the source file's `<title>` tag.
    - All other variables defined in config are available for substitution.

---

## 🚀 Usage

1.  Open a terminal in the `webclip-go` folder.
2.  Run the executable: `./webclip`
3.  Select your project from the list.
4.  Find your generated site in `projects/PROJECT_NAME/output/`.
