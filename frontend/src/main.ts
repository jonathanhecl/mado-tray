import "./style.css";
import {
  DisableStartup,
  EnableStartup,
  GetScripts,
  GetStartupStatus,
  RunScript,
  ToggleScript
} from "../wailsjs/go/main/App";
import type { main } from "../wailsjs/go/models";
import { WindowHide } from "../wailsjs/runtime/runtime";

type State = {
  scripts: main.Script[];
  startup: main.StartupStatus | null;
  loading: boolean;
  busyScriptId: string | null;
  startupBusy: boolean;
  error: string;
  notice: string;
};

const state: State = {
  scripts: [],
  startup: null,
  loading: true,
  busyScriptId: null,
  startupBusy: false,
  error: "",
  notice: ""
};

const root = document.querySelector<HTMLDivElement>("#app");

if (!root) {
  throw new Error("No se encontró el contenedor principal de la app.");
}

const appRoot = root;

async function load(): Promise<void> {
  state.loading = true;
  state.error = "";
  render();

  try {
    const [scripts, startup] = await Promise.all([GetScripts(), GetStartupStatus()]);
    state.scripts = scripts;
    state.startup = startup;
  } catch (error) {
    state.error = errorMessage(error);
  } finally {
    state.loading = false;
    render();
  }
}

function render(): void {
  appRoot.innerHTML = `
    <main class="panel">
      <header class="header">
        <div>
          <p class="eyebrow">Gestor de arranque macOS</p>
          <h1>Mado-Tray</h1>
        </div>
        <button class="icon-button" data-action="hide" title="Ocultar ventana">×</button>
      </header>

      ${state.error ? `<p class="message error">${escapeHtml(state.error)}</p>` : ""}
      ${state.notice ? `<p class="message notice">${escapeHtml(state.notice)}</p>` : ""}

      <section class="startup-card">
        <div>
          <h2>Abrir al iniciar macOS</h2>
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
      </section>

      <section class="scripts">
        <div class="section-title">
          <h2>Procesos</h2>
          <button class="ghost-button" data-action="reload" ${state.loading ? "disabled" : ""}>Recargar</button>
        </div>

        ${renderScripts()}
      </section>
    </main>
  `;
}

function renderScripts(): string {
  if (state.loading) {
    return `<p class="empty">Cargando configuración...</p>`;
  }

  if (state.scripts.length === 0) {
    return `<p class="empty">No hay procesos configurados todavía.</p>`;
  }

  return `
    <ul class="script-list">
      ${state.scripts.map(renderScript).join("")}
    </ul>
  `;
}

function renderScript(script: main.Script): string {
  const busy = state.busyScriptId === script.id;

  return `
    <li class="script-item">
      <div class="script-main">
        <div>
          <h3>${escapeHtml(script.name)}</h3>
          <p title="${escapeHtml(script.path)}">${escapeHtml(script.path)}</p>
        </div>
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
      <button class="run-button" data-action="run-script" data-id="${escapeHtml(script.id)}" ${busy ? "disabled" : ""}>
        ${busy ? "Abriendo..." : "Ejecutar ahora"}
      </button>
    </li>
  `;
}

function startupDescription(): string {
  if (!state.startup) {
    return "Revisando estado del sistema...";
  }

  if (!state.startup.available) {
    return escapeHtml(state.startup.message);
  }

  return state.startup.enabled
    ? "Mado-Tray ya está registrado como ítem de inicio."
    : "Activa este switch para registrar la app como ítem de inicio.";
}

appRoot.addEventListener("click", async (event) => {
  const target = event.target as HTMLElement;
  const button = target.closest<HTMLElement>("[data-action]");
  const action = button?.dataset.action;

  if (!action || action === "toggle-script" || action === "toggle-startup") {
    return;
  }

  if (action === "hide") {
    WindowHide();
    return;
  }

  if (action === "reload") {
    await load();
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
});

async function toggleScript(id: string, isActive: boolean): Promise<void> {
  state.busyScriptId = id;
  state.error = "";
  state.notice = "";
  render();

  try {
    state.scripts = await ToggleScript(id, isActive);
  } catch (error) {
    state.error = errorMessage(error);
    await load();
    return;
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
    await RunScript(id);
    state.notice = "Terminal.app abrió el proceso seleccionado.";
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
    state.startup = enabled ? await EnableStartup() : await DisableStartup();
    state.notice = enabled
      ? "Mado-Tray quedó configurado para abrir al iniciar macOS."
      : "Mado-Tray fue removido del arranque de macOS.";
  } catch (error) {
    state.error = errorMessage(error);
    state.startup = await GetStartupStatus();
  } finally {
    state.startupBusy = false;
    render();
  }
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
