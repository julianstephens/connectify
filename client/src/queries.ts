export const getHealthcheck = async (token: string) => {
  if (!token) throw new Error("auth token missing");
  const response = await fetch("/api/health", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  if (!response.ok) throw new Error("Healthcheck failed");
  return await response.json();
};
