# Mado-Tray

Mado-Tray es un gestor de procesos para el arranque de macOS. Vive en la barra superior, permite activar o desactivar scripts desde un panel pequeño y abre cada proceso en una ventana visible de Terminal.app para que puedas inspeccionarlo, detenerlo con `Ctrl+C` o interactuar con él manualmente.

Está construido con Go, Wails y un frontend TypeScript liviano.

## Funcionalidades

- Lee la configuración desde `~/.config/mado-tray/config.json`.
- Ejecuta automáticamente al iniciar la app todos los scripts con `is_active: true`.
- Abre los procesos en Terminal.app usando AppleScript, no como procesos ocultos.
- Permite activar/desactivar scripts y guardar el cambio en JSON.
- Permite crear y editar procesos desde un modal de la interfaz, además de eliminarlos.
- Incluye botón `Ejecutar ahora` para lanzar cualquier script manualmente.
- Soporta interfaz en español e inglés con selector `ES/EN`.
- Incluye switch `Abrir al iniciar macOS` para agregar o remover `Mado-Tray.app` de los ítems de inicio.
- Se muestra como panel flotante sin bordes y se controla desde el systray.
- Incluye `LSUIElement` en `build/darwin/Info.plist` para no aparecer en el Dock al empaquetar.

## Requisitos

- macOS.
- Go 1.23 o superior.
- Node.js 20 o superior recomendado.
- Wails CLI v2.

Instala Wails si no lo tienes:

```sh
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Verifica el entorno:

```sh
wails doctor
```

## Desarrollo

Instala dependencias:

```sh
npm install
go mod tidy
```

Ejecuta la app en modo desarrollo:

```sh
wails dev
```

En modo desarrollo, el switch `Abrir al iniciar macOS` puede mostrar un mensaje indicando que todavía no existe una `.app` final. Ese control está pensado para funcionar desde la app empaquetada.

## Configuración de procesos

Puedes administrar los procesos desde la interfaz:

1. Presiona `Agregar proceso`.
2. Completa `Nombre` y `Ruta` en el modal.
3. Activa `Activo al iniciar` si quieres que se ejecute cuando abra Mado-Tray.
4. Guarda el proceso.
5. Usa `Editar` o `Eliminar` en cada proceso cuando necesites cambiarlo.

La configuración también vive como JSON editable en:

```text
~/.config/mado-tray/config.json
```

Si el archivo no existe, Mado-Tray lo crea con un ejemplo inicial:

```json
[
  {
    "id": "example",
    "name": "Script de ejemplo",
    "path": "/ruta/a/tu/script.sh",
    "is_active": false
  }
]
```

Campos:

- `id`: identificador único del proceso.
- `name`: nombre visible en la interfaz.
- `path`: ruta absoluta al script o ejecutable.
- `is_active`: si es `true`, Mado-Tray lo abre automáticamente al iniciar.

Ejemplo real:

```json
[
  {
    "id": "api-local",
    "name": "API local",
    "path": "/Users/tu_usuario/Proyectos/api/start.sh",
    "is_active": true
  },
  {
    "id": "worker",
    "name": "Worker",
    "path": "/Users/tu_usuario/Proyectos/worker/run.sh",
    "is_active": false
  }
]
```

El script debe tener permisos de ejecución:

```sh
chmod +x /Users/tu_usuario/Proyectos/api/start.sh
```

## Idioma

La interfaz arranca en inglés por defecto. Si el sistema reporta español, Mado-Tray usa español automáticamente. También incluye un selector `ES/EN` en la parte superior; la preferencia se guarda localmente en el navegador embebido de Wails, así que se mantiene entre sesiones.

## Arranque con macOS

Desde la interfaz, usa el switch `Abrir al iniciar macOS`.

Cuando lo activas, Mado-Tray registra la `.app` actual como login item de macOS usando `System Events`. Cuando lo desactivas, elimina ese login item.

Para que funcione correctamente:

1. Compila la app.
2. Mueve `Mado-Tray.app` a `/Applications` o a una carpeta estable.
3. Abre esa app empaquetada.
4. Activa `Abrir al iniciar macOS`.

Si mueves la `.app` después de registrarla, desactiva y vuelve a activar el switch.

## Build

Genera la app:

```sh
wails build
```

El resultado queda en:

```text
build/bin/Mado-Tray.app
```

La plantilla `build/darwin/Info.plist` ya incluye:

```xml
<key>LSUIElement</key>
<true/>
```

Eso hace que Mado-Tray no aparezca en el Dock y viva principalmente en la barra superior.

## Publicación en GitHub

Checklist sugerido para publicar:

1. Ejecutar `npm install` y `go mod tidy`.
2. Ejecutar `npm run build`.
3. Ejecutar `wails build`.
4. Abrir `build/bin/Mado-Tray.app` y probar:
   - carga de `config.json`;
   - toggle de procesos;
   - botón `Ejecutar ahora`;
   - switch `Abrir al iniciar macOS`;
   - menú del systray.
5. Crear un tag, por ejemplo:

```sh
git tag v0.1.0
git push origin v0.1.0
```

6. Crear un release en GitHub y adjuntar la app empaquetada o un `.dmg` si decides distribuirla así.

## Estructura

```text
.
├── app.go                 # Métodos expuestos a Wails
├── config.go              # Lectura/escritura de config.json
├── runner.go              # Ejecución visible en Terminal.app
├── startup.go             # Login item de macOS
├── main.go                # Wails, ventana y systray
├── build/darwin/Info.plist
├── frontend/
│   ├── index.html
│   └── src/
│       ├── main.ts
│       └── style.css
└── wails.json
```

## Notas de seguridad

Mado-Tray ejecuta las rutas definidas en tu JSON. Usa rutas absolutas, revisa los scripts antes de activarlos y evita apuntar a archivos descargados de fuentes no confiables.