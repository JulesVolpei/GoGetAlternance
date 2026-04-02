import { MainLayout } from "@/layouts/MainLayout";
import { LogoCarousel } from "@/components/LogoCarousel";
import { SearchBar } from "@/components/SearchBar";
import type { SearchData } from "@/components/SearchBar";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import {useNavigate} from "react-router";

const Index = () => {
    const navigate = useNavigate(); // <-- 2. Initialisation du hook
    const searchMutation = useMutation({
        mutationFn: async (searchData: SearchData) => {
            const contractMapping: Record<string, string> = {
                "apprenticeship": "alternance",
                "internship": "stage"
            };

            const payload = {
                keywords: [searchData.query],
                contractTypes: [contractMapping[searchData.type] || searchData.type],
                platforms: searchData.sources
            };

            const response = await fetch("http://localhost:8080/api/cherche", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(payload),
            });

            if (!response.ok) {
                throw new Error("Erreur lors de la communication avec le serveur.");
            }

            return response.json();
        },
        onSuccess: (data) => {
            toast.success("Recherches terminées !", {
                description: `${data?.length || 0} offres ont été récupérées.`,
            });
            console.log("Offres récupérées :", data);
            navigate("/search", { state: { offers: data } });
        },
        onError: (error) => {
            toast.error("Oups, un problème est survenu.", {
                description: error.message,
            });
        }
    });

    const handleSearch = (data: SearchData) => {
        if (!data.query.trim()) {
            toast.warning("Veuillez saisir un mot-clé !");
            return;
        }

        toast.info("Lancement des scrapers en cours...", {
            description: `Recherche de "${data.query}" sur ${data.sources.length} plateformes.`
        });

        searchMutation.mutate(data);
    };

    return (
        <MainLayout>
            <div className="flex flex-col items-center justify-center min-h-[calc(100vh-3.5rem)] px-6 gap-12">
                <div className="text-center space-y-4 max-w-2xl">
                    <h1 className="text-4xl md:text-5xl font-bold tracking-tight">
                        GoGet<span className="text-muted-foreground">Alternance</span>
                    </h1>
                    <p className="text-muted-foreground text-lg">
                        Toutes les dernières offres d'alternances et de stages au même endroit.
                    </p>
                </div>

                <SearchBar onSearch={handleSearch} />

                {searchMutation.isPending && (
                    <div className="text-primary animate-pulse font-medium">
                        Recherche en cours... Veuillez patienter.
                    </div>
                )}

                <div className="w-full space-y-4">
                    <p className="text-center text-sm text-muted-foreground">
                        En provenance de
                    </p>
                    <LogoCarousel />
                </div>
            </div>
        </MainLayout>
    );
};

export default Index;