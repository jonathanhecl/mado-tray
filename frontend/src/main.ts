import "./style.css";
import type { main } from "../wailsjs/go/models";

type Locale = "es" | "en";
type BackendMethod =
  | "AddScript"
  | "DeleteScript"
  | "DisableStartup"
  | "EnableStartup"
  | "GetScripts"
  | "GetStartupStatus"
  | "HideWindow"
  | "PickScriptPath"
  | "RunScript"
  | "SetLocale"
  | "ToggleScript"
  | "UpdateScript";
type ScriptForm = {
  name: string;
  path: string;
  args: string;
  is_active: boolean;
};

type State = {
  scripts: main.Script[];
  startup: main.StartupStatus | null;
  locale: Locale;
  form: ScriptForm;
  editingId: string | null;
  deleteCandidateId: string | null;
  isFormOpen: boolean;
  isOptionsOpen: boolean;
  loading: boolean;
  busyScriptId: string | null;
  formBusy: boolean;
  startupBusy: boolean;
  error: string;
  notice: string;
};

type Dictionary = Record<string, string>;

const dictionaries: Record<Locale, Dictionary> = {
  es: {
    addProcess: "Agregar proceso",
    appConfigured: "Mado-Tray quedó configurado para abrir al iniciar macOS.",
    appRemoved: "Mado-Tray fue removido del arranque de macOS.",
    backendUnavailable: "El backend de Wails todavía no está disponible.",
    cancel: "Cancelar",
    acknowledge: "Aceptar",
    confirmDelete: "¿Eliminar este proceso de Mado-Tray?",
    confirmDeleteTitle: "Eliminar proceso",
    delete: "Eliminar",
    edit: "Editar",
    editingProcess: "Editando proceso",
    empty: "No hay procesos configurados todavía.",
    enabledAtStartup: "Mado-Tray ya está registrado como ítem de inicio.",
    hideWindow: "Ocultar ventana",
    language: "Idioma",
    loading: "Cargando configuración...",
    macStartup: "Abrir al iniciar macOS",
    manualRunNotice: "Terminal.app abrió el proceso seleccionado.",
    name: "Nombre",
    namePlaceholder: "API local",
    options: "Opciones",
    path: "Ruta",
    pathPlaceholder: "/Users/tu_usuario/Proyectos/api/start.sh",
    args: "Argumentos",
    argsPlaceholder: "--port 3000",
    browsePath: "Explorar",
    processAdded: "Proceso agregado.",
    processDeleted: "Proceso eliminado.",
    processes: "Procesos",
    processSaved: "Proceso guardado.",
    reload: "Recargar",
    reviewingStartup: "Revisando estado del sistema...",
    runNow: "Ejecutar ahora",
    running: "Abriendo...",
    saveChanges: "Guardar cambios",
    startupDisabled: "Activa este switch para registrar la app como ítem de inicio.",
    startupManager: "Gestor de arranque macOS",
    startWithMac: "Activo al iniciar",
    validationName: "Ingresa un nombre para el proceso.",
    validationPath: "Ingresa una ruta absoluta al script o ejecutable."
  },
  en: {
    addProcess: "Add process",
    appConfigured: "Mado-Tray is configured to open when macOS starts.",
    appRemoved: "Mado-Tray was removed from macOS startup.",
    backendUnavailable: "The Wails backend is not available yet.",
    cancel: "Cancel",
    acknowledge: "OK",
    confirmDelete: "Delete this process from Mado-Tray?",
    confirmDeleteTitle: "Delete process",
    delete: "Delete",
    edit: "Edit",
    editingProcess: "Editing process",
    empty: "No processes configured yet.",
    enabledAtStartup: "Mado-Tray is already registered as a login item.",
    hideWindow: "Hide window",
    language: "Language",
    loading: "Loading configuration...",
    macStartup: "Open at macOS startup",
    manualRunNotice: "Terminal.app opened the selected process.",
    name: "Name",
    namePlaceholder: "Local API",
    options: "Options",
    path: "Path",
    pathPlaceholder: "/Users/your_user/Projects/api/start.sh",
    args: "Arguments",
    argsPlaceholder: "--port 3000",
    browsePath: "Browse",
    processAdded: "Process added.",
    processDeleted: "Process deleted.",
    processes: "Processes",
    processSaved: "Process saved.",
    reload: "Reload",
    reviewingStartup: "Checking system status...",
    runNow: "Run now",
    running: "Opening...",
    saveChanges: "Save changes",
    startupDisabled: "Enable this switch to register the app as a login item.",
    startupManager: "macOS startup manager",
    startWithMac: "Active on startup",
    validationName: "Enter a process name.",
    validationPath: "Enter an absolute path to the script or executable."
  }
};

const state: State = {
  scripts: [],
  startup: null,
  locale: loadLocale(),
  form: emptyForm(),
  editingId: null,
  deleteCandidateId: null,
  isFormOpen: false,
  isOptionsOpen: false,
  loading: true,
  busyScriptId: null,
  formBusy: false,
  startupBusy: false,
  error: "",
  notice: ""
};

const root = document.querySelector<HTMLDivElement>("#app");

if (!root) {
  throw new Error("No se encontró el contenedor principal de la app.");
}

const appRoot = root;

let noticeTimer: ReturnType<typeof setTimeout> | null = null;
let armedNotice = "";

async function load(): Promise<void> {
  state.loading = true;
  state.error = "";
  render();

  try {
    await callBackend<void>("SetLocale", state.locale);
  } catch {
    // El backend puede no estar listo durante el primer render.
  }

  try {
    const [scripts, startup] = await Promise.all([
      callBackend<main.Script[]>("GetScripts"),
      callBackend<main.StartupStatus>("GetStartupStatus")
    ]);
    state.scripts = Array.isArray(scripts) ? scripts : [];
    state.startup = startup;
  } catch (error) {
    if (isBackendUnavailable(error)) {
      state.scripts = [];
      state.startup = unavailableStartupStatus();
    } else {
      state.error = errorMessage(error);
    }
  } finally {
    state.loading = false;
    render();
  }
}

function render(): void {
  appRoot.innerHTML = `
    <main class="panel">
      <header class="titlebar">
        <div class="titlebar-drag">
          <span class="titlebar-title">MadoTray</span>
        </div>
        <div class="titlebar-controls">
          <button class="titlebar-button titlebar-button-settings" data-action="open-options" title="${t("options")}" aria-label="${t("options")}">⚙</button>
          <button class="titlebar-button titlebar-button-close" data-action="hide" title="${t("hideWindow")}" aria-label="${t("hideWindow")}">×</button>
        </div>
      </header>

      <div class="panel-body">
        <section class="scripts">
          <div class="section-title">
            <h2>${t("processes")}</h2>
            <div class="section-actions">
              <button class="primary-button" data-action="open-add-process">${t("addProcess")}</button>
              <button class="icon-action-button" data-action="reload" title="${t("reload")}" aria-label="${t("reload")}" ${state.loading ? "disabled" : ""}>🔄</button>
            </div>
          </div>

          ${renderScripts()}
        </section>
      </div>

      ${state.isFormOpen ? renderFormModal() : ""}
      ${state.isOptionsOpen ? renderOptionsModal() : ""}
      ${state.deleteCandidateId ? renderDeleteModal() : ""}
      ${renderToasts()}
    </main>
  `;

  armNoticeDismiss();
}

function renderToasts(): string {
  if (state.error) {
    return `
      <div class="toast-host">
        <div class="toast toast-error" role="alert">
          <span class="toast-message">${escapeHtml(state.error)}</span>
          <button class="toast-dismiss" type="button" data-action="dismiss-error">${t("acknowledge")}</button>
        </div>
      </div>
    `;
  }

  if (!state.notice) {
    return "";
  }

  return `
    <div class="toast-host">
      <div class="toast toast-success" role="status">
        <span class="toast-message">${escapeHtml(state.notice)}</span>
      </div>
    </div>
  `;
}

function clearNoticeTimer(): void {
  if (noticeTimer !== null) {
    clearTimeout(noticeTimer);
    noticeTimer = null;
  }
}

function armNoticeDismiss(): void {
  if (!state.notice || state.error) {
    armedNotice = "";
    clearNoticeTimer();
    return;
  }

  if (state.notice === armedNotice && noticeTimer !== null) {
    return;
  }

  armedNotice = state.notice;
  clearNoticeTimer();
  noticeTimer = setTimeout(() => {
    noticeTimer = null;
    armedNotice = "";
    state.notice = "";
    render();
  }, 3000);
}

function renderDeleteModal(): string {
  const script = state.scripts.find((item) => item.id === state.deleteCandidateId);
  const name = script?.name ?? "";

  return `
    <div class="modal-backdrop" data-action="close-delete">
      <section class="form-card modal" role="dialog" aria-modal="true" aria-labelledby="delete-title">
        <div class="modal-header">
          <h2 id="delete-title">${t("confirmDeleteTitle")}</h2>
          <button class="icon-button" type="button" data-action="close-delete" title="${t("cancel")}">×</button>
        </div>
        <p class="modal-copy">${t("confirmDelete")}</p>
        ${name ? `<p class="delete-target">${escapeHtml(name)}</p>` : ""}
        <div class="modal-actions">
          <button class="ghost-button" type="button" data-action="close-delete">${t("cancel")}</button>
          <button class="danger-button" type="button" data-action="confirm-delete" ${state.busyScriptId ? "disabled" : ""}>${t("delete")}</button>
        </div>
      </section>
    </div>
  `;
}

function renderOptionsModal(): string {
  return `
    <div class="modal-backdrop" data-action="close-options">
      <section class="form-card modal" role="dialog" aria-modal="true" aria-labelledby="options-title">
        <div class="modal-header">
          <h2 id="options-title">${t("options")}</h2>
          <button class="icon-button" type="button" data-action="close-options" title="${t("cancel")}">×</button>
        </div>

        <div class="settings-list">
          <div class="settings-row">
            <div>
              <h3>${t("language")}</h3>
            </div>
            <label class="language-select">
              <select data-action="change-locale" aria-label="${t("language")}">
                <option value="es" ${state.locale === "es" ? "selected" : ""}>ES</option>
                <option value="en" ${state.locale === "en" ? "selected" : ""}>EN</option>
              </select>
            </label>
          </div>

          <div class="settings-row">
            <div>
              <h3>${t("macStartup")}</h3>
              <p>${startupDescription()}</p>
            </div>
            <label class="switch">
              <input
                type="checkbox"
                data-action="toggle-startup"
                ${state.startup?.enabled ? "checked" : ""}
                ${state.startupBusy || !state.startup?.available ? "disabled" : ""}
              />
              <span></span>
            </label>
          </div>
        </div>
      </section>
    </div>
  `;
}

function renderFormModal(): string {
  const isEditing = state.editingId !== null;

  return `
    <div class="modal-backdrop" data-action="close-form">
      <section class="form-card modal" role="dialog" aria-modal="true" aria-labelledby="process-form-title">
        <div class="modal-header">
          <h2 id="process-form-title">${isEditing ? t("editingProcess") : t("addProcess")}</h2>
          <button class="icon-button" type="button" data-action="close-form" title="${t("cancel")}">×</button>
        </div>
        <form data-action="script-form">
          <label class="field">
            <span>${t("name")}</span>
            <input
              name="name"
              type="text"
              value="${escapeHtml(state.form.name)}"
              placeholder="${t("namePlaceholder")}"
              ${state.formBusy ? "disabled" : ""}
            />
          </label>

          <label class="field">
            <span>${t("path")}</span>
            <div class="path-field">
              <input
                name="path"
                type="text"
                value="${escapeHtml(state.form.path)}"
                placeholder="${t("pathPlaceholder")}"
                ${state.formBusy ? "disabled" : ""}
              />
              <button class="ghost-button path-browse-button" type="button" data-action="pick-path" ${state.formBusy ? "disabled" : ""}>
                ${t("browsePath")}
              </button>
            </div>
          </label>

          <label class="field">
            <span>${t("args")}</span>
            <input
              name="args"
              type="text"
              value="${escapeHtml(state.form.args)}"
              placeholder="${t("argsPlaceholder")}"
              ${state.formBusy ? "disabled" : ""}
            />
          </label>

          <div class="form-row">
            <label class="inline-check">
              <input name="is_active" type="checkbox" ${state.form.is_active ? "checked" : ""} ${state.formBusy ? "disabled" : ""} />
              <span>${t("startWithMac")}</span>
            </label>
            <div class="form-actions">
              <button class="ghost-button" type="button" data-action="close-form" ${state.formBusy ? "disabled" : ""}>${t("cancel")}</button>
              <button class="primary-button" type="submit" ${state.formBusy ? "disabled" : ""}>
                ${isEditing ? t("saveChanges") : t("addProcess")}
              </button>
            </div>
          </div>
        </form>
      </section>
    </div>
  `;
}

function renderScripts(): string {
  if (state.loading) {
    return `<p class="empty">${t("loading")}</p>`;
  }

  if (state.scripts.length === 0) {
    return `<p class="empty">${t("empty")}</p>`;
  }

  return `
    <ul class="script-list">
      ${state.scripts.map(renderScript).join("")}
    </ul>
  `;
}

function renderScript(script: main.Script): string {
  const busy = state.busyScriptId === script.id;
  const path = script.path?.trim() ?? "";
  const args = script.args?.trim() ?? "";
  const commandTitle = args ? `${path} ${args}` : path;

  return `
    <li class="script-item">
      <div class="script-header">
        <h3>${escapeHtml(script.name)}</h3>
        <label class="switch">
          <input
            type="checkbox"
            data-action="toggle-script"
            data-id="${escapeHtml(script.id)}"
            ${script.is_active ? "checked" : ""}
            ${busy ? "disabled" : ""}
          />
          <span></span>
        </label>
      </div>
      <div class="script-meta">
        <p class="script-line script-path" title="${escapeHtml(path)}">
          <span class="script-line-label">${t("path")}</span>
          <span class="script-line-value">${escapeHtml(path)}</span>
        </p>
        ${args ? `
          <p class="script-line script-args" title="${escapeHtml(args)}">
            <span class="script-line-label">${t("args")}</span>
            <span class="script-line-value">${escapeHtml(args)}</span>
          </p>
        ` : ""}
      </div>
      <div class="script-footer">
        <div class="script-actions">
          <button class="ghost-button" data-action="edit-script" data-id="${escapeHtml(script.id)}" ${busy ? "disabled" : ""}>${t("edit")}</button>
          <button class="danger-button" data-action="delete-script" data-id="${escapeHtml(script.id)}" ${busy ? "disabled" : ""}>${t("delete")}</button>
        </div>
        <button class="run-button" data-action="run-script" data-id="${escapeHtml(script.id)}" title="${escapeHtml(commandTitle)}" ${busy ? "disabled" : ""}>
          ${busy ? t("running") : t("runNow")}
        </button>
      </div>
    </li>
  `;
}

function startupDescription(): string {
  if (!state.startup) {
    return t("reviewingStartup");
  }

  if (!state.startup.available) {
    return escapeHtml(state.startup.message);
  }

  return state.startup.enabled
    ? t("enabledAtStartup")
    : t("startupDisabled");
}

appRoot.addEventListener("click", async (event) => {
  const target = event.target as HTMLElement;
  const button = target.closest<HTMLElement>("[data-action]");
  const action = button?.dataset.action;

  if (!action || action === "toggle-script" || action === "toggle-startup") {
    return;
  }

  if (action === "hide") {
    try {
      await callBackend<void>("HideWindow");
    } catch {
      window.runtime?.WindowHide?.();
    }
    return;
  }

  if (action === "reload") {
    await load();
    return;
  }

  if (action === "dismiss-error") {
    state.error = "";
    render();
    return;
  }

  if (action === "open-options") {
    state.isOptionsOpen = true;
    state.error = "";
    state.notice = "";
    render();
    return;
  }

  if (action === "close-options") {
    if (button.classList.contains("modal-backdrop") && target.closest(".modal")) {
      return;
    }
    state.isOptionsOpen = false;
    render();
    return;
  }

  if (action === "open-add-process") {
    state.isFormOpen = true;
    state.editingId = null;
    state.form = emptyForm();
    state.error = "";
    state.notice = "";
    render();
    return;
  }

  if (action === "pick-path") {
    syncFormFromDom();
    await pickScriptPath();
    return;
  }

  if (action === "close-form") {
    if (button.classList.contains("modal-backdrop") && target.closest(".modal")) {
      return;
    }
    resetForm();
    render();
    return;
  }

  if (action === "edit-script") {
    const id = button.dataset.id;
    if (id) {
      startEditing(id);
    }
    return;
  }

  if (action === "delete-script") {
    const id = button.dataset.id;
    if (id) {
      state.deleteCandidateId = id;
      state.error = "";
      state.notice = "";
      render();
    }
    return;
  }

  if (action === "close-delete") {
    if (button.classList.contains("modal-backdrop") && target.closest(".modal")) {
      return;
    }
    state.deleteCandidateId = null;
    render();
    return;
  }

  if (action === "confirm-delete") {
    if (state.deleteCandidateId) {
      await deleteScript(state.deleteCandidateId);
    }
    return;
  }

  if (action === "run-script") {
    const id = button.dataset.id;
    if (id) {
      await runScript(id);
    }
  }
});

appRoot.addEventListener("change", async (event) => {
  const input = event.target as HTMLInputElement;
  const action = input.dataset.action;

  if (action === "toggle-script") {
    const id = input.dataset.id;
    if (id) {
      await toggleScript(id, input.checked);
    }
  }

  if (action === "toggle-startup") {
    await toggleStartup(input.checked);
  }

  if (action === "change-locale") {
    const locale = input.value === "en" ? "en" : "es";
    state.locale = locale;
    localStorage.setItem("mado-tray-locale", locale);
    try {
      await callBackend<void>("SetLocale", locale);
    } catch {
      // Si falla el backend, mantenemos al menos el idioma de la UI.
    }
    render();
  }
});

appRoot.addEventListener("submit", async (event) => {
  const form = event.target as HTMLFormElement;
  if (form.dataset.action !== "script-form") {
    return;
  }

  event.preventDefault();
  await submitScriptForm(form);
});

async function pickScriptPath(): Promise<void> {
  try {
    const picked = await callBackend<string>("PickScriptPath");
    if (!picked) {
      return;
    }

    state.form.path = picked;
    render();

    requestAnimationFrame(() => {
      const input = appRoot.querySelector<HTMLInputElement>('input[name="args"]');
      if (!input) {
        return;
      }

      input.focus();
      input.setSelectionRange(input.value.length, input.value.length);
    });
  } catch (error) {
    state.error = errorMessage(error);
    render();
  }
}

function syncFormFromDom(): void {
  const form = appRoot.querySelector<HTMLFormElement>('form[data-action="script-form"]');
  if (!form) {
    return;
  }

  const nameInput = form.querySelector<HTMLInputElement>('input[name="name"]');
  const pathInput = form.querySelector<HTMLInputElement>('input[name="path"]');
  const argsInput = form.querySelector<HTMLInputElement>('input[name="args"]');
  const activeInput = form.querySelector<HTMLInputElement>('input[name="is_active"]');

  if (nameInput) {
    state.form.name = nameInput.value;
  }
  if (pathInput) {
    state.form.path = pathInput.value;
  }
  if (argsInput) {
    state.form.args = argsInput.value;
  }
  if (activeInput) {
    state.form.is_active = activeInput.checked;
  }
}

async function toggleScript(id: string, isActive: boolean): Promise<void> {
  state.busyScriptId = id;
  state.error = "";
  state.notice = "";
  render();

  try {
    state.scripts = await callBackend<main.Script[]>("ToggleScript", id, isActive);
  } catch (error) {
    state.error = errorMessage(error);
    await load();
    return;
  } finally {
    state.busyScriptId = null;
    render();
  }
}

async function submitScriptForm(form: HTMLFormElement): Promise<void> {
  const formData = new FormData(form);
  const input = {
    name: String(formData.get("name") ?? "").trim(),
    path: String(formData.get("path") ?? "").trim(),
    args: String(formData.get("args") ?? "").trim(),
    is_active: formData.get("is_active") === "on"
  };
  state.form = input;

  if (!input.name) {
    state.error = t("validationName");
    state.notice = "";
    render();
    return;
  }

  if (!input.path) {
    state.error = t("validationPath");
    state.notice = "";
    render();
    return;
  }

  state.formBusy = true;
  state.error = "";
  state.notice = "";
  render();

  try {
    state.scripts = state.editingId
      ? await callBackend<main.Script[]>("UpdateScript", state.editingId, input)
      : await callBackend<main.Script[]>("AddScript", input);
    state.notice = state.editingId ? t("processSaved") : t("processAdded");
    resetForm();
  } catch (error) {
    state.error = errorMessage(error);
  } finally {
    state.formBusy = false;
    render();
  }
}

function startEditing(id: string): void {
  const script = state.scripts.find((item) => item.id === id);
  if (!script) {
    return;
  }

  state.editingId = id;
  state.isFormOpen = true;
  state.form = {
    name: script.name,
    path: script.path,
    args: script.args ?? "",
    is_active: script.is_active
  };
  state.error = "";
  state.notice = "";
  render();
}

async function deleteScript(id: string): Promise<void> {
  state.busyScriptId = id;
  state.error = "";
  state.notice = "";
  render();

  try {
    state.scripts = await callBackend<main.Script[]>("DeleteScript", id);
    if (state.editingId === id) {
      resetForm();
    }
    state.deleteCandidateId = null;
    state.notice = t("processDeleted");
  } catch (error) {
    state.error = errorMessage(error);
  } finally {
    state.busyScriptId = null;
    render();
  }
}

async function runScript(id: string): Promise<void> {
  state.busyScriptId = id;
  state.error = "";
  state.notice = "";
  render();

  try {
    await callBackend<void>("RunScript", id);
    state.notice = t("manualRunNotice");
  } catch (error) {
    state.error = errorMessage(error);
  } finally {
    state.busyScriptId = null;
    render();
  }
}

async function toggleStartup(enabled: boolean): Promise<void> {
  state.startupBusy = true;
  state.error = "";
  state.notice = "";
  render();

  try {
    state.startup = enabled
      ? await callBackend<main.StartupStatus>("EnableStartup")
      : await callBackend<main.StartupStatus>("DisableStartup");
    state.notice = enabled ? t("appConfigured") : t("appRemoved");
  } catch (error) {
    state.error = errorMessage(error);
    state.startup = await callBackend<main.StartupStatus>("GetStartupStatus");
  } finally {
    state.startupBusy = false;
    render();
  }
}

function resetForm(): void {
  state.editingId = null;
  state.isFormOpen = false;
  state.form = emptyForm();
}

function emptyForm(): ScriptForm {
  return {
    name: "",
    path: "",
    args: "",
    is_active: false
  };
}

function loadLocale(): Locale {
  const savedLocale = localStorage.getItem("mado-tray-locale");
  if (savedLocale === "es" || savedLocale === "en") {
    return savedLocale;
  }

  return navigator.language.toLowerCase().startsWith("es") ? "es" : "en";
}

function t(key: string): string {
  return dictionaries[state.locale][key] ?? key;
}

async function callBackend<T>(method: BackendMethod, ...args: unknown[]): Promise<T> {
  const app = window.go?.main?.App;
  const fn = app?.[method];

  if (!fn) {
    throw new Error(t("backendUnavailable"));
  }

  return (await fn(...args)) as T;
}

function isBackendUnavailable(error: unknown): boolean {
  return error instanceof Error && error.message === t("backendUnavailable");
}

function unavailableStartupStatus(): main.StartupStatus {
  return {
    enabled: false,
    app_path: "",
    available: false,
    message: t("reviewingStartup")
  } as main.StartupStatus;
}

function errorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }

  return String(error);
}

function escapeHtml(value: string): string {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

void load();
