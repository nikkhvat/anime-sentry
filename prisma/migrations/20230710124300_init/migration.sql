-- CreateTable
CREATE TABLE "Anime" (
    "id" SERIAL NOT NULL,
    "url" TEXT NOT NULL,
    "lastEpisode" INTEGER NOT NULL DEFAULT 0,
    "releaseDate" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Anime_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "AnimeSubscription" (
    "id" SERIAL NOT NULL,
    "telegramUserId" TEXT NOT NULL,
    "animeId" INTEGER NOT NULL,

    CONSTRAINT "AnimeSubscription_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "Anime_url_key" ON "Anime"("url");

-- AddForeignKey
ALTER TABLE "AnimeSubscription" ADD CONSTRAINT "AnimeSubscription_animeId_fkey" FOREIGN KEY ("animeId") REFERENCES "Anime"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
