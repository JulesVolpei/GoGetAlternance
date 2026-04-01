import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";
import { AppSidebar } from "@/components/AppSidebar";

interface MainLayoutProps {
    children: React.ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {

    return (
        <SidebarProvider>
            <div className="min-h-screen flex w-full">
                <AppSidebar/>
                <div className="flex-1 flex flex-col">
                    <header className="h-14 flex items-center border-b px-4">
                        <SidebarTrigger />
                    </header>
                    <main className="flex-1">{children}</main>
                </div>
            </div>
        </SidebarProvider>
    );
}
