import { useState, useEffect } from "react";
import { useLocation, Link } from "react-router-dom";
import { MainLayout } from "@/layouts/MainLayout";
import { Building2, MapPin, ExternalLink, Briefcase, ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";

interface JobOffer {
    title: string;
    company: string;
    contract: string;
    location: string;
    url: string;
    scrapeDate: string;
    source: string;
}

const shuffleArray = (array: JobOffer[]) => {
    const shuffled = [...array];
    for (let i = shuffled.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
    }
    return shuffled;
};

export default function Search() {
    const location = useLocation();
    const initialOffers: JobOffer[] = location.state?.offers || [];

    const [shuffledOffers, setShuffledOffers] = useState<JobOffer[]>([]);
    const [currentPage, setCurrentPage] = useState(1);
    const itemsPerPage = 20;

    useEffect(() => {
        if (initialOffers.length > 0) {
            setShuffledOffers(shuffleArray(initialOffers));
            setCurrentPage(1); // On s'assure de revenir à la page 1
        }
    }, [initialOffers]);

    const totalPages = Math.ceil(shuffledOffers.length / itemsPerPage);
    const startIndex = (currentPage - 1) * itemsPerPage;
    const currentOffers = shuffledOffers.slice(startIndex, startIndex + itemsPerPage);

    const handlePageChange = (newPage: number) => {
        setCurrentPage(newPage);
        window.scrollTo({ top: 0, behavior: "smooth" });
    };

    return (
        <MainLayout>
            <div className="container mx-auto py-8 px-4 max-w-5xl">
                <div className="mb-8">
                    <h1 className="text-3xl font-bold tracking-tight mb-2">Résultats de recherche</h1>
                    <p className="text-muted-foreground">
                        {shuffledOffers.length} offre{shuffledOffers.length > 1 ? 's' : ''} trouvée{shuffledOffers.length > 1 ? 's' : ''}
                    </p>
                </div>

                {shuffledOffers.length === 0 ? (
                    <div className="text-center py-20 border rounded-xl bg-card">
                        <p className="text-lg text-muted-foreground mb-4">Aucune offre à afficher.</p>
                        <Button asChild>
                            <Link to="/">Retour à l'accueil</Link>
                        </Button>
                    </div>
                ) : (
                    <>
                        <div className="grid gap-4 md:grid-cols-2 mb-8">
                            {currentOffers.map((offer, index) => (
                                <div key={index} className="flex flex-col p-6 rounded-xl border bg-card hover:shadow-md transition-shadow">
                                    <div className="flex-1">
                                        <div className="flex justify-between items-start gap-4 mb-4">
                                            <h2 className="text-xl font-semibold line-clamp-2">{offer.title}</h2>
                                            <span className="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold text-muted-foreground whitespace-nowrap bg-secondary/50">
                                                {offer.source}
                                            </span>
                                        </div>

                                        <div className="space-y-2 text-sm text-muted-foreground mb-6">
                                            <div className="flex items-center gap-2">
                                                <Building2 className="h-4 w-4 text-primary/70" />
                                                <span className="font-medium text-foreground">{offer.company}</span>
                                            </div>
                                            <div className="flex items-center gap-2">
                                                <MapPin className="h-4 w-4 text-primary/70" />
                                                <span>{offer.location}</span>
                                            </div>
                                            <div className="flex items-center gap-2">
                                                <Briefcase className="h-4 w-4 text-primary/70" />
                                                <span>{offer.contract}</span>
                                            </div>
                                        </div>
                                    </div>

                                    <div className="pt-4 border-t mt-auto">
                                        <Button asChild className="w-full" variant="default">
                                            <a href={offer.url} target="_blank" rel="noopener noreferrer">
                                                Voir l'offre
                                                <ExternalLink className="h-4 w-4 ml-2" />
                                            </a>
                                        </Button>
                                    </div>
                                </div>
                            ))}
                        </div>

                        {totalPages > 1 && (
                            <div className="flex items-center justify-center gap-4 py-4">
                                <Button
                                    variant="outline"
                                    onClick={() => handlePageChange(currentPage - 1)}
                                    disabled={currentPage === 1}
                                >
                                    <ChevronLeft className="h-4 w-4 mr-2" />
                                    Précédent
                                </Button>

                                <span className="text-sm font-medium text-muted-foreground">
                                    Page {currentPage} sur {totalPages}
                                </span>

                                <Button
                                    variant="outline"
                                    onClick={() => handlePageChange(currentPage + 1)}
                                    disabled={currentPage === totalPages}
                                >
                                    Suivant
                                    <ChevronRight className="h-4 w-4 ml-2" />
                                </Button>
                            </div>
                        )}
                    </>
                )}
            </div>
        </MainLayout>
    );
}