import {
    Carousel,
    CarouselContent,
    CarouselItem,
    CarouselNext,
    CarouselPrevious,
} from "./ui/carousel.tsx";

const platforms = [
    { name: "Indeed", logo: "https://upload.wikimedia.org/wikipedia/commons/f/fc/Indeed_logo.svg" },
    { name: "Welcome to the Jungle", logo: "https://cdn.welcometothejungle.com/wttj-front/production/assets/images/logos/wttj.svg" },
    { name: "HelloWork", logo: "https://upload.wikimedia.org/wikipedia/commons/e/e1/LOGO_HelloWork_Activit%C3%A9s.png" },
    { name: "FranceTravail", logo: "https://upload.wikimedia.org/wikipedia/fr/thumb/8/8d/France-travail-2023.svg/1280px-France-travail-2023.svg.png" },
    // { name: "LinkedIn", logo: "https://upload.wikimedia.org/wikipedia/commons/0/01/LinkedIn_Logo.svg" },
    // { name: "Glassdoor", logo: "https://upload.wikimedia.org/wikipedia/commons/e/e1/Glassdoor_logo.svg" },
];

export function LogoCarousel() {
    return (
        <div className="w-full max-w-2xl mx-auto">
            <Carousel
                opts={{ align: "center" }}
                className="w-full"
            >
                <CarouselContent className="-ml-4">
                    {platforms.map((platform) => (
                        <CarouselItem key={platform.name} className="pl-4 basis-1/3 flex items-center justify-center">
                            <div className="flex flex-col items-center gap-3 p-6 rounded-lg border bg-card h-28 w-full justify-center">
                                <img
                                    src={platform.logo}
                                    alt={platform.name}
                                    className="h-8 max-w-[120px] object-contain grayscale hover:grayscale-0 transition-all"
                                />
                                <span className="text-xs text-muted-foreground">{platform.name}</span>
                            </div>
                        </CarouselItem>
                    ))}
                </CarouselContent>
                <CarouselPrevious />
                <CarouselNext />
            </Carousel>
        </div>
    );
}
