import { useState } from 'react'
import { createRoot } from 'react-dom/client'
import { CargaDocumental } from './componentes/carga-documental/carga-documental'
import { Documentos } from './componentes/documentos/documentos'
import { iniciarSesion } from './desarrollo/sesion.api'
import './style.css'
function App() {
  const [usuario, setUsuario] = useState(''),
    [revision, setRevision] = useState(0),
    [error, setError] = useState(''),
    [ocupado, setOcupado] = useState(false)
  async function cambiar(valor: string) {
    setOcupado(true)
    setError('')
    try {
      await iniciarSesion(valor)
      setUsuario(valor)
      setRevision((r) => r + 1)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Sesión no disponible')
    } finally {
      setOcupado(false)
    }
  }
  return (
    <>
      <header>
        <a className="marca" href="/">
          <span>▤</span> documental<span className="punto">.</span>
        </a>
        <div className="entorno">
          GRUPO 07 <span>DESARROLLO LOCAL</span>
        </div>
      </header>
      <main>
        <div className="intro">
          <div>
            <div className="etiqueta">BÓVEDA DOCUMENTAL</div>
            <h1>
              Documentos a salvo.
              <br />
              <span>Confianza verificable.</span>
            </h1>
            <p>
              Un espacio para conservar tus archivos originales,
              <br className="solo-escritorio" /> con una firma que protege su integridad.
            </p>
          </div>
          <div className="sesion">
            <label htmlFor="usuario">Identidad de prueba</label>
            <select
              id="usuario"
              value={usuario}
              disabled={ocupado}
              onChange={(e) => void cambiar(e.target.value)}
            >
              <option value="" disabled>
                Selecciona una identidad
              </option>
              <option value="operador">Operador · propietario</option>
              <option value="lector">Lector · mismo propietario</option>
              <option value="ajeno">Operador · otro propietario</option>
            </select>
            <small>Proveedor local de desarrollo</small>
          </div>
        </div>
        {error && (
          <p role="alert" className="error">
            {error}
          </p>
        )}
        {usuario ? (
          <div className="contenido" key={usuario}>
            <CargaDocumental
              baseUrl=""
              puedeCargar={usuario !== 'lector'}
              alCargar={() => setRevision((r) => r + 1)}
            />
            <Documentos baseUrl="" revision={revision} />
          </div>
        ) : (
          <section className="panel bienvenida">
            <h2>Comienza con una identidad de prueba</h2>
            <p>
              Selecciona un usuario para explorar la carga, consulta y verificación de
              documentos.
            </p>
          </section>
        )}
        <div className="garantias">
          <span>◇ Originales preservados</span>
          <span>◇ Firma criptográfica</span>
          <span>◇ Acceso controlado</span>
        </div>
      </main>
      <footer>
        MS Gestor Documental <span>Base común LPS · PINGESO 2026</span>
      </footer>
    </>
  )
}
createRoot(document.getElementById('root')!).render(<App />)
