import { PageLayout } from "@/components/layout";
import LandingPage from "@/pages/Landing";
import NotFoundPage from "@/pages/NotFound";
import UserProfilePage from "@/pages/UserProfile";
import { Route, Routes } from "react-router-dom";

export const AppRoutes = () => (
  <Routes>
    <Route index path="/" element={<LandingPage />} />
    <Route
      path="/settings"
      element={
        <PageLayout>
          <UserProfilePage />
        </PageLayout>
      }
    />
    <Route path="*" element={<NotFoundPage />} />
  </Routes>
);
