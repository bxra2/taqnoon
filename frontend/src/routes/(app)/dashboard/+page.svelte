<script lang="ts">
    import { BookOpen, Languages, Database } from 'lucide-svelte'

    let { data } = $props()

    let stats = $derived([
        {
            label: 'إجمالي المصطلحات',
            value: (data.stats?.totalTerms ?? 0).toLocaleString('ar-EG'),
            icon: Database,
            color: 'text-blue-500',
            bg: 'bg-blue-500/10'
        },
        {
            label: 'المصطلحات عربية',
            value: (data.stats?.arabic ?? 0).toLocaleString('ar-EG'),
            icon: Languages,
            color: 'text-amber-500',
            bg: 'bg-amber-500/10'
        },
        {
            label: 'المصطلحات فرنسية',
            value: (data.stats?.french ?? 0).toLocaleString('ar-EG'),
            icon: Languages,
            color: 'text-blue-500',
            bg: 'bg-blue-500/10'
        },
        {
            label: 'المصطلحات الألمانية',
            value: (data.stats?.german ?? 0).toLocaleString('ar-EG'),
            icon: Languages,
            color: 'text-purple-500',
            bg: 'bg-purple-500/10'
        },  
        {
            label: 'المعاجم النشطة',
            value: (data.stats?.glossaries ?? 0).toLocaleString('ar-EG'),
            icon: BookOpen,
            color: 'text-emerald-500',
            bg: 'bg-emerald-500/10'
        },
    ])
</script>

<div class="flex flex-col gap-8 animate-in fade-in slide-in-from-bottom-4 duration-500 relative left-1/2 right-1/2 -ml-[-5vw] -mr-[40vw]">
    <!-- Welcome Header -->
    <div class="flex flex-col gap-2">
        <h1 class="text-3xl font-bold almahdi">لوحة المعلومات</h1>
        <p class="text-gray-500 dark:text-gray-400">نظرة عامة على محتوى المعاجم والمصطلحات</p>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {#each stats as stat}
            <div class="p-6 rounded-2xl border accent shadow-sm hover:shadow-md">
                <div class="flex items-start justify-between">
                    <div class="flex flex-col gap-3">
                        <span class="text-gray-500 dark:text-gray-400 text-sm font-medium">{stat.label}</span>
                        <span class="text-3xl font-bold text-gray-900 dark:text-white dir-ltr font-['Inter']">
                            {stat.value}
                        </span>
                    </div>
                    <div class={`p-3 rounded-xl ${stat.bg} ${stat.color}`}>
                        <stat.icon size={24} />
                    </div>
                </div>
            </div>
        {/each}
    </div>

    <!-- Additional Content Placeholder (can be expanded later) -->
    <!-- Using CSS Grid for layout flexibility -->
    <!-- <div class="grid grid-cols-1 lg:grid-cols-3 gap-6"> -->
        <!-- Main Chart Area (Placeholder) -->
        <!-- <div class="lg:col-span-2 p-6 rounded-2xl border border-gray-100 dark:border-stone-800 bg-white dark:bg-stone-900 min-h-[300px] flex items-center justify-center relative overflow-hidden group">
             <div class="absolute inset-0 bg-gradient-to-br from-gray-50 to-transparent dark:from-stone-800/50 opacity-50"></div>
             <div class="text-center relative z-10">
                <Activity class="w-12 h-12 mx-auto text-gray-300 mb-4 group-hover:scale-110 transition-transform duration-500" />
                <p class="text-gray-400 font-medium">إحصائيات البحث قريباً</p>
             </div>
        </div> -->

        <!-- Side Card (e.g. Tips or Updates) -->
        <!-- <div class="p-6 rounded-2xl border border-gray-100 dark:border-stone-800 bg-gradient-to-br from-primary/5 to-transparent dark:from-primary/10">
            <h3 class="font-bold text-lg mb-4">هل تعلم؟</h3>
            <p class="text-gray-600 dark:text-gray-300 leading-relaxed text-sm">
                يمكنك استخدام علامة التنصيص "..." للبحث عن عبارات دقيقة، أو استخدام عوامل التصفية لتضييق نطاق البحث حسب المعجم.
            </p>
        </div>
    </div> -->
</div>
