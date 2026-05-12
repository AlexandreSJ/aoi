<div align="center">
  <img width=100% src="https://capsule-render.vercel.app/api?type=waving&height=200&color=0:00aaff,100:00aaff&text=A%20O%20I&fontColor=ffffff&fontSize=50&fontAlignY=40" alt="AOI">
</div>

<h1 align="center">🔹 あおい 🔹</h1>

<p align="center">
  Un test de mecanografía en la terminal.
  <br>
  Practica tu escritura, relájate y disfruta con aoi.
</p>

<div align="center">

  <img src="https://img.shields.io/badge/Go-1.24-00add8?style=flat-square&logo=go" alt="Go" href="https://go.dev/doc/install">
  <img src="https://img.shields.io/badge/Bubble_Tea-1.3-ff69b4?style=flat-square" alt="Bubble Tea" href="https://github.com/charmbracelet/bubbletea">
  <img src="https://img.shields.io/badge/Lipgloss-1.1-7d56f4?style=flat-square" alt="Lipgloss" href="https://github.com/charmbracelet/lipgloss">
  <img src="https://img.shields.io/badge/Licencia-MIT-00aaff?style=flat-square" alt="Licencia" href="LICENSE">
  <br>
  <a href="https://www.buymeacoffee.com/aelxand" target="_blank">
    <img width=120px src="assets/bmc/bmc.png" alt="Buy Me A Coffee">
  </a>
  <br>
  <a href="README-ptBR.md"><img src="https://img.shields.io/badge/ptBR🇧🇷-README-fff?style=flat-square" alt="License" href="LICENSE"></a>
  <a href="README.md"><img src="https://img.shields.io/badge/en🇬🇧-README-fff?style=flat-square" alt="License" href="LICENSE"></a>
  <a href="README-es.md"><img src="https://img.shields.io/badge/es🇪🇸-README-0af?style=flat-square" alt="License" href="LICENSE"></a>
</div>

## Índice

- [¿Qué es Aoi?](#qué-es-aoi)
- [Instalación](#instalación)
- [Uso](#uso)
- [Características](#características)
- [Requisitos del Sistema](#requisitos-del-sistema)
- [Solución de Problemas](#solución-de-problemas)
- [Hoja de Ruta](#hoja-de-ruta)
- [Contribuir](#contribuir)
- [Licencia](#licencia)

## ¿Qué es Aoi?

Empecé a disfrutar haciendo tests de mecanografía como hobby y para mantener mis habilidades de escritura afiladas, pero siempre quise esto en una TUI. ¡Así que creé AOI!

Elige entre 4 modos diferentes de práctica de mecanografía en Aoi:
- Zen: Escribe infinitamente a tu propio ritmo
- Temporizado: Corre contra el reloj
- Contado: Escribe un número fijo de palabras
- Cita: Escribe una cita aleatoria

Configura los colores como quieras. También puedes agregar más palabras o citas, ¡escalable para usar cualquier idioma!

<div align="center">
  <img width=100% src="assets/prints/typing.png" alt="Escribiendo">
</div>

## Instalación

### Requisitos Previos

- Go 1.24+ (necesario para compilar desde el código fuente)
- Emulador de terminal con soporte para colores ANSI y Unicode

### Métodos de Instalación

#### Método 1: Instalar usando Go

```bash
# Asegúrate de tener Go configurado en tu ~/.zshrc o ~/.bashrc
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Cierra la terminal o ejecuta para aplicar
source ~/.zshrc
# O
source ~/.bashrc

# Finalmente, instala directamente desde GitHub
go install -a github.com/AlexandreSJ/aoi/cmd/aoi@latest
```

#### Método 2: Compilar desde el Código Fuente (para desarrolladores)

```bash
# Clona el repositorio
git clone https://github.com/AlexandreSJ/aoi.git
cd aoi

# Compila la aplicación
make build
```

### Inicio Rápido

Después de la instalación, simplemente ejecuta:

```bash
aoi
```

### Comandos de Build

Si tienes el repositorio git instalado, en /aoi puedes ejecutar:

```bash
make clean  # Eliminar el directorio /build
make build  # Compilar el binario
make run    # Compilar y ejecutar inmediatamente
```

### Características

- **Ligero y rápido** - Ligero y veloz como un erizo
- **Feedback de escritura en tiempo real** - Ve tu precisión y velocidad mientras escribes
- **Soporte Unicode** - Funciona con varios conjuntos de caracteres
- **Diseño responsivo** - Se adapta a diferentes tamaños de terminal

### Requisitos del Sistema

- **Sistema Operativo**: Linux, macOS o Windows (con WSL)
- **Terminal**: Cualquier emulador de terminal moderno (Terminal, iTerm2, Alacritty, Windows Terminal, etc.)
- **Espacio en Disco**: ~5MB para el binario

### Solución de Problemas

**P: Recibo "command not found: aoi"**
R: Asegúrate de que el directorio GOPATH/bin esté en tu PATH, o usa la ruta completa al binario.

**P: Los colores se ven extraños en mi terminal**
R: Intenta configurar `TERM=xterm-256color` o usa una terminal con soporte true color.

**P: La aplicación no inicia**
R: Asegúrate de tener Go 1.24+ instalado y que tu terminal soporta caracteres Unicode.

**P: Tengo problemas para instalar/actualizar aoi a la última versión**
R: Si ya tienes Go instalado, ejecuta el siguiente comando para evitar proxy.golang.org y usar la etiqueta `-a` para forzar la recompilación:
`GOPROXY=direct go install -a github.com/AlexandreSJ/aoi/cmd/aoi@latest`

## Hoja de Ruta

> Plan de desarrollo por fases. Cada fase mezcla trabajo de múltiples áreas.

### Fase 1 — Persistencia de Configuración y Temas `v0.3`


**Gestión de Configuración** (Completado para v0.3.1)
- Mejoras en el padding y espaciado de la UI de entrada
- Correcciones en la entrada de configuración

**Persistencia de configuración** (En progreso)
- Recordar último modo usado, duración temporizada, conteo de palabras, último archivo por modo, preferencia de tamaño de texto
- Campo de versión en config con migración automática en cambios de schema
- Opción cursor activado/desactivado
- Opción flash de tecla activado/desactivado
- Refactorizar manejo de errores en `config.go`

**Sistema de temas**
- 4+ temas predefinidos (Aoi predeterminado, Monokai, Dracula, Solarized, Catppuccin)
- Pantalla de selección de tema accesible desde home o configuración
- Vista previa del tema antes de aplicar
- Temas personalizados desde `~/.config/aoi/themes/` como archivos YAML
- Importar/exportar archivos de tema para compartir

### Fase 2 — Telemetría y Estadísticas `v0.4`

**Almacenamiento de telemetría (JSONL)**
- Una línea por sesión de prueba en `~/.config/aoi/telemetry.jsonl`
- Por sesión: modo, WPM, caracteres escritos, palabras completas, ok/errores, precisión, duración, timestamp
- Timestamps por carácter para análisis de velocidad
- Solo adición, sin dependencias externas

**Pantalla de estadísticas**
- Nueva pantalla accesible desde home (tecla `s`)
- Promedio de WPM por modo
- Lista de historial de sesiones — navegable por fecha/modo/WPM
- Mapa de calor por tecla — teclas más lentas y con más errores
- Mejores resultados por modo — mostrados en la pantalla home
- Seguimiento de racha — días consecutivos de práctica, mostrados en home
- Exportar telemetría como CSV

**Gráficos y visualizaciones**
- Gráfico de WPM a lo largo del tiempo (sparkline o gráfico de barras)
- Tendencia de precisión
- Barras de rendimiento por tecla
- Todo renderizado con Lipgloss/Bubble Tea (sin lib de gráficos externa)

### Fase 3 — UX/UI y Sonido `v0.5`

**Barra de progreso**
- Barra de progreso estilizada para modo Temporizado (tiempo restante) y modo Contado (palabras escritas/restantes)
- Renderizada arriba o abajo del área de escritura

**Refactorización de la pantalla de escritura**
- Experiencia de lectura más suave — mejor espaciado entre líneas, renderizado de espacio entre palabras
- Indicador de WPM en tiempo real durante la escritura
- Refactorizar visualización del timer (modo Temporizado) para mostrar en el área de la barra de progreso en vez del pie de página
- Repetición de palabras con error — después del test, opción de repracticar solo palabras con error

**Sistema de sonido**
- `gopxl/beep` para reproducción de audio
- 3+ paquetes de sonidos integrados (toque mecánico, clic suave, etc.)
- Controles separados: sonido de escritura ok, sonido de escritura error, jingle de finalización
- Sonidos personalizados desde `~/.config/aoi/sounds/`

**Refactorización del pie de página**
- Mover información de tiempo/modo a otras ubicaciones (área de la barra de progreso)
- Pie de página más limpio con solo atajos de teclado

## Contribuir

¡Las contribuciones son bienvenidas! No dudes en enviar un Pull Request.

## Licencia

MIT

<div align="center">
  <a href="https://git.io/typing-svg">
    <img src="https://readme-typing-svg.herokuapp.com?font=Fira+Code&duration=1&color=00AAFF&center=true&vCenter=true&repeat=false&width=435&lines=stay+blue+%3C3" alt="Typing SVG" />
  </a>
</div>

<img width=100% src="https://capsule-render.vercel.app/api?type=slice&height=300&color=00aaff&text=AOI&section=footer&fontAlign=22&fontAlignY=69&rotate=19&fontSize=50&fontColor=ffffff&desc=あおい&descAlignY=80&descAlign=22" alt="AOI">
