import { createFileRoute } from '@tanstack/react-router'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { AccountTab } from './-components/account-tab'
import { OccupationTab } from './-components/occupation-tab'
import { ProfileTab } from './-components/profile-tab'

export const Route = createFileRoute('/asesi/profil/')({
  component: ProfilPage,
})

function ProfilPage() {
  const { auth } = Route.useRouteContext()

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-lg font-semibold text-slate-900">Profil Asesi</h1>
        <p className="text-sm text-slate-500">Kelola data akun, profil, dan pekerjaan Anda.</p>
      </div>

      <Tabs defaultValue="account">
        <TabsList>
          <TabsTrigger value="account">Data Akun</TabsTrigger>
          <TabsTrigger value="profile">Data Profil</TabsTrigger>
          <TabsTrigger value="occupation">Pekerjaan</TabsTrigger>
        </TabsList>

        <TabsContent value="account">
          <AccountTab session={auth} />
        </TabsContent>
        <TabsContent value="profile">
          <ProfileTab session={auth} />
        </TabsContent>
        <TabsContent value="occupation">
          <OccupationTab session={auth} />
        </TabsContent>
      </Tabs>
    </div>
  )
}
