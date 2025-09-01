package eshelper

func EnsureMap(m map[string]any, key string) map[string]any {
   if _, exists := m[key]; !exists {
      m[key] = make(map[string]any)
   }
   return m[key].(map[string]any)
}
