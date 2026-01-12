// Gerador de ID simples para uso temporário no frontend
// Os IDs reais são gerados pela API

export const generateId = (): string => {
  const timestamp = Date.now().toString(36);
  const randomPart = Math.random().toString(36).substring(2, 15);
  return `${timestamp}-${randomPart}`;
};
