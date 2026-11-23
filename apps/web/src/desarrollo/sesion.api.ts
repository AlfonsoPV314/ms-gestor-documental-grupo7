export async function iniciarSesion(usuario: string) {
  const r = await fetch('/identidad/sesion', {
    method: 'POST',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ usuario }),
  })
  if (!r.ok) throw Error('No se pudo iniciar la sesión de desarrollo')
}
