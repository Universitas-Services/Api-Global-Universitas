import sys
import json
import urllib.request

def search_skills(query=None):
    url = "https://api.github.com/repos/sickn33/antigravity-awesome-skills/contents/skills"
    
    req = urllib.request.Request(url, headers={'User-Agent': 'Antigravity-Agent-IDE'})
    
    try:
        with urllib.request.urlopen(req) as response:
            items = json.loads(response.read().decode())
            
            if not isinstance(items, list):
                print("Error: La API de GitHub no devolvió una lista de elementos.")
                return

            # Filtrar por elementos que sean directorios (carpetas de skills)
            skills = [item for item in items if item.get('type') == 'dir']
            
            # Si se proporciona una búsqueda, filtrar por el nombre
            if query:
                query_lower = query.lower()
                skills = [s for s in skills if query_lower in s.get('name', '').lower()]
            
            if not skills:
                mensaje = f"No se encontraron skills para: {query}" if query else "No se encontraron carpetas de skills en el repositorio."
                print(mensaje)
                return

            print("--- Resultados de Skills en sickn33/antigravity-awesome-skills ---")
            for item in skills:
                print(f"\nNombre: {item.get('name')}")
                print(f"URL: {item.get('html_url')}")
                
    except Exception as e:
        print(f"Error al consultar el repositorio en GitHub: {str(e)}")

if __name__ == "__main__":
    search_term = sys.argv[1] if len(sys.argv) > 1 else None
    search_skills(search_term)
