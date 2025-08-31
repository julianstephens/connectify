import { useAuth0 } from "@auth0/auth0-react";
import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { getHealthcheck } from "./queries";

export const useAuthToken = () => {
  const [token, setToken] = useState("");
  const { getAccessTokenSilently, getAccessTokenWithPopup } = useAuth0();

  useEffect(() => {
    if (!token) {
      (async () => {
        try {
          const t = await getAccessTokenSilently({
            authorizationParams: {
              audience: import.meta.env.VITE_AUTH0_AUDIENCE,
            },
          });
          setToken(t);
        } catch (e) {
          console.error(e);
          try {
            const t = await getAccessTokenWithPopup({
              authorizationParams: {
                audience: import.meta.env.VITE_AUTH0_AUDIENCE,
              },
            });
            if (t) {
              setToken(t);
            }
          } catch (err) {
            console.error(e);
          }
        }
      })();
    }
  }, [getAccessTokenSilently, getAccessTokenWithPopup, token]);

  return token;
};

export const useHealthcheck = () => {
  const token = useAuthToken();
  return useQuery({
    queryKey: ["healthcheck"],
    queryFn: () => getHealthcheck(token),
    enabled: !!token,
  });
};
